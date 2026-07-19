package training

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrQRExpired  = errors.New("QR Code sudah expired")
var ErrAlreadyScan = errors.New("Anda sudah melakukan absensi pada sesi ini")
var ErrNotFound   = errors.New("data tidak ditemukan")

type Repository interface {
	ListSesi(ctx context.Context, unitID string, tanggal string) ([]Sesi, error)
	CreateSesi(ctx context.Context, req CreateSesiRequest) (string, error)
	UpdateSesi(ctx context.Context, id string, req CreateSesiRequest) error
	DeleteSesi(ctx context.Context, id string) error
	GenerateQR(ctx context.Context, sesiID string, durasiJam int) (*QRAbsensi, error)
	ScanQR(ctx context.Context, req ScanQRRequest) error
	GetKehadiran(ctx context.Context, sesiID string) ([]Kehadiran, error)
	RingkasanKehadiran(ctx context.Context, sesiID string, totalAnggota int) (*RingkasanKehadiran, error)
}

type pgRepository struct{ db *sql.DB }
func NewRepository(db *sql.DB) Repository { return &pgRepository{db: db} }

func (r *pgRepository) ListSesi(ctx context.Context, unitID, tanggal string) ([]Sesi, error) {
	where := "WHERE 1=1"
	args := []any{}
	i := 1
	if unitID != "" { where += fmt.Sprintf(" AND s.unit_id=$%d", i); args = append(args, unitID); i++ }
	if tanggal != "" { where += fmt.Sprintf(" AND s.tanggal=$%d::date", i); args = append(args, tanggal); i++ }
	_ = i
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT s.id, s.unit_id, u.nama, u.cabang_id::text, COALESCE(s.pelatih_id::text,''), COALESCE(a.nama_lengkap,''), 
		       s.jenis, s.tanggal::text, s.jam_mulai::text, COALESCE(s.jam_selesai::text,''), 
		       COALESCE(s.lokasi,''), COALESCE(s.materi,''), s.created_at
		FROM sesi_latihan s 
		JOIN unit_latihan u ON u.id = s.unit_id 
		LEFT JOIN pelatih p ON p.id = s.pelatih_id
		LEFT JOIN anggota a ON a.id = p.anggota_id
		%s ORDER BY s.tanggal DESC, s.jam_mulai DESC
	`, where), args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []Sesi
	for rows.Next() {
		var s Sesi
		if err := rows.Scan(&s.ID, &s.UnitID, &s.UnitNama, &s.CabangID, &s.PelatihID, &s.PelatihNama, &s.Jenis, &s.Tanggal, &s.JamMulai, &s.JamSelesai, &s.Lokasi, &s.Materi, &s.CreatedAt); err != nil { return nil, err }
		result = append(result, s)
	}
	return result, nil
}

func (r *pgRepository) CreateSesi(ctx context.Context, req CreateSesiRequest) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO sesi_latihan (unit_id, pelatih_id, jenis, tanggal, jam_mulai, jam_selesai, lokasi, materi)
		VALUES ($1, $2, $3, $4::date, $5::time, NULLIF($6,'')::time, $7, $8) RETURNING id
	`, req.UnitID, req.PelatihID, req.Jenis, req.Tanggal, req.JamMulai, req.JamSelesai, req.Lokasi, req.Materi).Scan(&id)
	return id, err
}

func (r *pgRepository) UpdateSesi(ctx context.Context, id string, req CreateSesiRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sesi_latihan 
		SET unit_id=$1, pelatih_id=$2, jenis=$3, tanggal=$4::date, jam_mulai=$5::time, jam_selesai=NULLIF($6,'')::time, lokasi=$7, materi=$8
		WHERE id=$9
	`, req.UnitID, req.PelatihID, req.Jenis, req.Tanggal, req.JamMulai, req.JamSelesai, req.Lokasi, req.Materi, id)
	return err
}

func (r *pgRepository) DeleteSesi(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sesi_latihan WHERE id=$1", id)
	return err
}

func (r *pgRepository) GenerateQR(ctx context.Context, sesiID string, durasiJam int) (*QRAbsensi, error) {
	token := uuid.New().String()
	exp := time.Now().Add(time.Duration(durasiJam) * time.Hour)
	// Deactivate old QR for this session first
	_, _ = r.db.ExecContext(ctx, "UPDATE absensi_qr SET is_active=FALSE WHERE sesi_id=$1", sesiID)
	var qr QRAbsensi
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO absensi_qr (sesi_id, qr_token, expires_at) VALUES ($1, $2, $3)
		RETURNING id, sesi_id, qr_token, expires_at, is_active
	`, sesiID, token, exp).Scan(&qr.ID, &qr.SesiID, &qr.QRToken, &qr.ExpiresAt, &qr.IsActive)
	return &qr, err
}

func (r *pgRepository) ScanQR(ctx context.Context, req ScanQRRequest) error {
	var expiresAt time.Time
	var isActive bool
	var sesiID string
	err := r.db.QueryRowContext(ctx, `SELECT sesi_id, expires_at, is_active FROM absensi_qr WHERE qr_token=$1`, req.QRToken).Scan(&sesiID, &expiresAt, &isActive)
	if errors.Is(err, sql.ErrNoRows) { return ErrNotFound }
	if !isActive || time.Now().After(expiresAt) { return ErrQRExpired }
	// Check duplicate
	var exists bool
	_ = r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM kehadiran WHERE sesi_id=$1 AND anggota_id=$2)`, sesiID, req.AnggotaID).Scan(&exists)
	if exists { return ErrAlreadyScan }
	_, err = r.db.ExecContext(ctx, `INSERT INTO kehadiran (sesi_id, anggota_id, metode) VALUES ($1, $2, 'qr')`, sesiID, req.AnggotaID)
	return err
}

func (r *pgRepository) GetKehadiran(ctx context.Context, sesiID string) ([]Kehadiran, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT k.id, k.sesi_id, k.anggota_id, a.nama_lengkap, k.waktu_scan, k.metode
		FROM kehadiran k JOIN anggota a ON a.id = k.anggota_id
		WHERE k.sesi_id=$1 ORDER BY k.waktu_scan ASC
	`, sesiID)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []Kehadiran
	for rows.Next() {
		var k Kehadiran
		if err := rows.Scan(&k.ID, &k.SesiID, &k.AnggotaID, &k.Nama, &k.WaktuScan, &k.Metode); err != nil { return nil, err }
		result = append(result, k)
	}
	return result, nil
}

func (r *pgRepository) RingkasanKehadiran(ctx context.Context, sesiID string, totalAnggota int) (*RingkasanKehadiran, error) {
	var hadir int
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM kehadiran WHERE sesi_id=$1`, sesiID).Scan(&hadir)
	pct := 0.0
	if totalAnggota > 0 { pct = float64(hadir) / float64(totalAnggota) * 100 }
	return &RingkasanKehadiran{SesiID: sesiID, Hadir: hadir, Tidak: totalAnggota - hadir, Persentase: pct}, nil
}
