// Package training manages Jadwal Latihan and QR Absensi.
package training

import "time"

type Sesi struct {
	ID          string    `json:"id"`
	UnitID      string    `json:"unit_id"`
	UnitNama    string    `json:"unit_nama,omitempty"`
	CabangID    string    `json:"cabang_id,omitempty"`
	PelatihID   string    `json:"pelatih_id"`
	PelatihNama string    `json:"pelatih_nama,omitempty"`
	Jenis       string    `json:"jenis"` // Latihan Rutin, Khusus, Pelatih
	Tanggal     string    `json:"tanggal"`
	JamMulai    string    `json:"jam_mulai"`
	JamSelesai  string    `json:"jam_selesai"`
	Lokasi      string    `json:"lokasi"`
	Materi      string    `json:"materi,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateSesiRequest struct {
	UnitID    string `json:"unit_id"`
	PelatihID string `json:"pelatih_id"`
	Jenis     string `json:"jenis"`
	Tanggal   string `json:"tanggal"`
	JamMulai  string `json:"jam_mulai"`
	JamSelesai string `json:"jam_selesai"`
	Lokasi    string `json:"lokasi"`
	Materi    string `json:"materi"`
}

type QRAbsensi struct {
	ID        string    `json:"id"`
	SesiID    string    `json:"sesi_id"`
	QRToken   string    `json:"qr_token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

type ScanQRRequest struct {
	QRToken   string `json:"qr_token"`
	AnggotaID string `json:"anggota_id"`
}

type Kehadiran struct {
	ID        string    `json:"id"`
	SesiID    string    `json:"sesi_id"`
	AnggotaID string    `json:"anggota_id"`
	Nama      string    `json:"nama,omitempty"`
	WaktuScan time.Time `json:"waktu_scan"`
	Metode    string    `json:"metode"`
}

type RingkasanKehadiran struct {
	SesiID    string `json:"sesi_id"`
	Hadir     int    `json:"hadir"`
	Tidak     int    `json:"tidak"`
	Persentase float64 `json:"persentase"`
}
