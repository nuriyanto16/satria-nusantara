// Package organization manages Cabang, Unit, Anggota, and Pelatih data.
package organization

import "time"

// ─── Cabang ──────────────────────────────────────────────────────────────────

type Cabang struct {
	ID         string    `json:"id"`
	Kode       string    `json:"kode"`
	Nama       string    `json:"nama"`
	ProvinsiID int       `json:"provinsi_id"`
	KotaKab    string    `json:"kota_kab"`
	Alamat     string    `json:"alamat,omitempty"`
	Status     string    `json:"status"`
	// Aggregates (not stored, computed on query)
	JumlahUnit    int `json:"jumlah_unit,omitempty"`
	JumlahAnggota int `json:"jumlah_anggota,omitempty"`
	JumlahPelatih int `json:"jumlah_pelatih,omitempty"`
}

type CreateCabangRequest struct {
	Kode       string `json:"kode"`
	Nama       string `json:"nama"`
	ProvinsiID int    `json:"provinsi_id"`
	KotaKab    string `json:"kota_kab"`
	Alamat     string `json:"alamat"`
}

// ─── Unit Latihan ──────────────────────────────────────────────────────────────

type Unit struct {
	ID        string `json:"id"`
	CabangID  string `json:"cabang_id"`
	Nama      string `json:"nama"`
	LokasiNama string `json:"lokasi_nama"`
	LokasiAlamat string `json:"lokasi_alamat"`
	PicUserID string `json:"pic_user_id,omitempty"`
	Status    string `json:"status"`
	// Aggregates
	JumlahAnggota int `json:"jumlah_anggota,omitempty"`
}

type CreateUnitRequest struct {
	CabangID     string `json:"cabang_id"`
	Nama         string `json:"nama"`
	LokasiNama   string `json:"lokasi_nama"`
	LokasiAlamat string `json:"lokasi_alamat"`
	PicUserID    string `json:"pic_user_id"`
}

// ─── Anggota ──────────────────────────────────────────────────────────────────

type TingkatanEnum string

const (
	TingPraDasar  TingkatanEnum = "Pra Dasar"
	TingDasar     TingkatanEnum = "Dasar"
	TingPH        TingkatanEnum = "PH"
	TingGabungan  TingkatanEnum = "Gabungan"
	TingPK        TingkatanEnum = "PK"
	TingGPK       TingkatanEnum = "GPK"
	TingPenjuru   TingkatanEnum = "Penjuru"
	TingSelangkah TingkatanEnum = "Selangkah"
	TingMeditasi  TingkatanEnum = "Meditasi"
)

type Anggota struct {
	ID              string        `json:"id"`
	UserID          *string       `json:"user_id,omitempty"`
	NomorAnggota    string        `json:"nomor_anggota"`
	NamaLengkap     string        `json:"nama_lengkap"`
	TanggalLahir    *string       `json:"tanggal_lahir,omitempty"`
	JenisKelamin    string        `json:"jenis_kelamin"`
	NoHp            string        `json:"no_hp"`
	FotoURL         string        `json:"foto_url,omitempty"`
	UnitID          string        `json:"unit_id"`
	UnitNama        string        `json:"unit_nama,omitempty"`   // joined
	CabangNama      string        `json:"cabang_nama,omitempty"` // joined
	Tingkatan       TingkatanEnum `json:"tingkatan"`
	JurusSaatIni   int           `json:"jurus_saat_ini"`
	CounterKehadiran int          `json:"counter_kehadiran"`
	Status          string        `json:"status"`
	TanggalDaftar   time.Time     `json:"tanggal_daftar"`
	TanggalAktif    *time.Time    `json:"tanggal_aktif,omitempty"`
}

type CreateAnggotaRequest struct {
	NamaLengkap  string `json:"nama_lengkap"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	NoHp         string `json:"no_hp"`
	UnitID       string `json:"unit_id"`
}

type VerifikasiAnggotaRequest struct {
	Aksi string `json:"aksi"` // "approve" | "reject"
}

// ─── Pelatih ─────────────────────────────────────────────────────────────────

type Pelatih struct {
	ID                string   `json:"id"`
	AnggotaID         string   `json:"anggota_id"`
	NamaLengkap       string   `json:"nama_lengkap"` // joined from anggota
	NomorAnggota      string   `json:"nomor_anggota,omitempty"` // joined from anggota
	Tingkatan         string   `json:"tingkatan,omitempty"`     // joined from anggota
	CabangID          string   `json:"cabang_id"`
	Jenis             string   `json:"jenis"` // cabang | pusat
	KategoriTransport string   `json:"kategori_transport"`
	IsActive          bool     `json:"is_active"`
	UnitList          []string `json:"unit_list,omitempty"` // unit IDs
}

// ─── Pagination ───────────────────────────────────────────────────────────────

type ListParams struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	CabangID string `json:"cabang_id"`
	UnitID   string `json:"unit_id"`
	Status   string `json:"status"`
}

type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// ─── Sebaran Provinsi ──────────────────────────────────────────────────────────

type SebaranProvinsi struct {
	ID            int    `json:"id"`
	Provinsi      string `json:"provinsi"`
	CabangUtama   string `json:"cabang_utama"`
	JumlahUnit    int    `json:"jumlah_unit"`
	JumlahAnggota int    `json:"jumlah_anggota"`
}

// ─── Dashboard Stats ──────────────────────────────────────────────────────────

type DashboardStats struct {
	TotalCabang          int                  `json:"total_cabang"`
	TotalAnggota         int                  `json:"total_anggota"`
	TotalUnit            int                  `json:"total_unit"`
	NewAnggotaThisWeek   int                  `json:"new_anggota_this_week"`
	NewCabangThisMonth   int                  `json:"new_cabang_this_month"`
	AttendanceRate       int                  `json:"attendance_rate"`
	StatGender           DashboardGender      `json:"stat_gender"`
	StatUsia             DashboardUsia        `json:"stat_usia"`
	RekapBLBA            []DashboardBLBA      `json:"rekap_blba"`
	TopUnits             []DashboardTopUnit   `json:"top_units"`
	RankingKebugaran     []DashboardRanking   `json:"ranking_kebugaran"`
	TrenKebugaran        DashboardTren        `json:"tren_kebugaran"`
}

type DashboardGender struct {
	LakiLaki  int `json:"laki_laki"`
	Perempuan int `json:"perempuan"`
}

type DashboardUsia struct {
	Under17   int `json:"under_17"`
	Usia18_25 int `json:"usia_18_25"`
	Usia26_40 int `json:"usia_26_40"`
	Over40    int `json:"over_40"`
}

type DashboardBLBA struct {
	Cabang    string `json:"cabang"`
	Terkumpul int64  `json:"terkumpul"`
	Target    int64  `json:"target"`
	Pct       int    `json:"pct"`
	Class     string `json:"class"`
}

type DashboardTopUnit struct {
	Nama   string `json:"nama"`
	Cabang string `json:"cabang"`
	Pct    int    `json:"pct"`
}

type DashboardRanking struct {
	Cabang   string `json:"cabang"`
	Skor     int    `json:"skor"`
	Kategori string `json:"kategori"`
}

type DashboardTren struct {
	Labels   []string           `json:"labels"`
	Datasets []DashboardDataset `json:"datasets"`
}

type DashboardDataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

type AnggotaStats struct {
	BLBABulanIni              string               `json:"blba_bulan_ini"`
	BLBALunas                 int                  `json:"blba_lunas"`
	BLBATotal                 int                  `json:"blba_total"`
	KehadiranBulanIni         int                  `json:"kehadiran_bulan_ini"`
	KehadiranTotalBulanIni    int                  `json:"kehadiran_total_bulan_ini"`
	TotalKehadiran            int                  `json:"total_kehadiran"`
	KehadiranTarget           int                  `json:"kehadiran_target"`
	KehadiranPct              int                  `json:"kehadiran_pct"`
	AttendanceHistory         []AttendanceMonth    `json:"attendance_history"`
	Score                     int                  `json:"score"`
	Category                  string               `json:"category"`
	CatClass                  string               `json:"cat_class"`
	DiffStr                   string               `json:"diff_str"`
	DiffColor                 string               `json:"diff_color"`
	Trends                    []int                `json:"trends"`
	Points                    []Point              `json:"points"`
	PolylinePoints            string               `json:"polyline_points"`
	PolygonPoints             string               `json:"polygon_points"`
	Results                   []FitnessResult      `json:"results"`
}

type AttendanceMonth struct {
	Month    string `json:"month"`
	Attended int    `json:"attended"`
	Total    int    `json:"total"`
	Pct      int    `json:"pct"`
}

type Point struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Score int     `json:"score"`
}

type FitnessResult struct {
	Label string    `json:"label"`
	Value string    `json:"value"`
	Cat   ResultCat `json:"cat"`
}

type ResultCat struct {
	Label string `json:"label"`
	Class string `json:"class"`
}



type UpdateKebugaranRequest struct {
	NafasDalamAir        int     `json:"nafas_dalam_air"`
	PushUp              int     `json:"push_up"`
	SitUp               int     `json:"sit_up"`
	SitAndReach         int     `json:"sit_and_reach"`
	ShuttleRun          float64 `json:"shuttle_run"`
	KategoriKeseluruhan string  `json:"kategori_keseluruhan"`
}
