package organization

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("data tidak ditemukan")

// ErrDuplicateKode is returned when creating a cabang with an existing kode.
var ErrDuplicateKode = errors.New("kode cabang sudah digunakan")

// Service defines the business logic for the organization domain.
type Service interface {
	// Cabang
	ListCabang(params ListParams, userScope, userCabangID string) (*PaginatedResult[Cabang], error)
	GetCabang(id string) (*Cabang, error)
	GetCabangTrends(id string, period int) (*CabangTrendsResponse, error)
	CreateCabang(req CreateCabangRequest, createdBy string) (string, error)
	UpdateCabang(id string, req CreateCabangRequest) error

	// Unit
	ListUnit(cabangID string) ([]Unit, error)
	GetUnit(id string) (*Unit, error)
	CreateUnit(req CreateUnitRequest) (string, error)

	// Anggota
	ListAnggota(params ListParams) (*PaginatedResult[Anggota], error)
	GetAnggota(id string) (*Anggota, error)
	CreateAnggota(req CreateAnggotaRequest) (string, error)
	UpdateAnggota(id string, req UpdateAnggotaRequest) error
	VerifikasiAnggota(id string, req VerifikasiAnggotaRequest) error
	UpdateFotoAnggota(id, fotoURL string) error
	GetAnggotaStats(id string) (*AnggotaStats, error)
	UpdateAnggotaKebugaran(id string, req UpdateKebugaranRequest) error
	GetAnggotaKebugaranHistory(id string) ([]map[string]interface{}, error)

	// Pelatih
	ListPelatih(cabangID string) ([]Pelatih, error)

	// Sebaran
	GetSebaranProvinsi() ([]SebaranProvinsi, error)

	GetDashboardStats() (*DashboardStats, error)

	UpdateUnit(id string, req CreateUnitRequest) error
	DeleteUnit(id string) error
	
	ListPengurus(cabangID string) ([]PengurusCabang, error)
	CreatePengurus(req CreatePengurusRequest) (string, error)
	UpdatePengurus(id string, req CreatePengurusRequest) error
	DeletePengurus(id string) error
	
	CreatePelatih(req CreatePelatihRequest) (string, error)
	UpdatePelatih(id string, req CreatePelatihRequest) error
	DeletePelatih(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service { return &service{repo: repo} }

// ─── Cabang ─────────────────────────────────────────────────────────────────

func (s *service) ListCabang(params ListParams, userScope, userCabangID string) (*PaginatedResult[Cabang], error) {
	// Pengurus Cabang can only see their own cabang
	if userScope == "cabang" {
		params.CabangID = userCabangID
	}
	return s.repo.ListCabang(context.Background(), params)
}

func (s *service) GetCabang(id string) (*Cabang, error) {
	return s.repo.GetCabangByID(context.Background(), id)
}

func (s *service) GetCabangTrends(id string, period int) (*CabangTrendsResponse, error) {
	return s.repo.GetCabangTrends(context.Background(), id, period)
}


func (s *service) CreateCabang(req CreateCabangRequest, createdBy string) (string, error) {
	return s.repo.CreateCabang(context.Background(), req, createdBy)
}

func (s *service) UpdateCabang(id string, req CreateCabangRequest) error {
	return s.repo.UpdateCabang(context.Background(), id, req)
}

// ─── Unit ────────────────────────────────────────────────────────────────────

func (s *service) ListUnit(cabangID string) ([]Unit, error) {
	return s.repo.ListUnit(context.Background(), cabangID)
}

func (s *service) GetUnit(id string) (*Unit, error) {
	return s.repo.GetUnitByID(context.Background(), id)
}

func (s *service) CreateUnit(req CreateUnitRequest) (string, error) {
	return s.repo.CreateUnit(context.Background(), req)
}

func (s *service) UpdateUnit(id string, req CreateUnitRequest) error {
	return s.repo.UpdateUnit(context.Background(), id, req)
}

func (s *service) DeleteUnit(id string) error {
	return s.repo.DeleteUnit(context.Background(), id)
}

// ─── Anggota ─────────────────────────────────────────────────────────────────

func (s *service) ListAnggota(params ListParams) (*PaginatedResult[Anggota], error) {
	return s.repo.ListAnggota(context.Background(), params)
}

func (s *service) GetAnggota(id string) (*Anggota, error) {
	return s.repo.GetAnggotaByID(context.Background(), id)
}

func (s *service) CreateAnggota(req CreateAnggotaRequest) (string, error) {
	return s.repo.CreateAnggota(context.Background(), req)
}

func (s *service) UpdateAnggota(id string, req UpdateAnggotaRequest) error {
	return s.repo.UpdateAnggota(context.Background(), id, req)
}

func (s *service) VerifikasiAnggota(id string, req VerifikasiAnggotaRequest) error {
	if req.Aksi != "approve" && req.Aksi != "reject" && req.Aksi != "aktifkan" && req.Aksi != "nonaktifkan" {
		return errors.New("aksi harus 'approve', 'reject', 'aktifkan', atau 'nonaktifkan'")
	}
	status := "aktif"
	if req.Aksi == "reject" || req.Aksi == "nonaktifkan" {
		status = "nonaktif"
	}
	return s.repo.VerifikasiAnggota(context.Background(), id, status)
}

func (s *service) UpdateFotoAnggota(id, fotoURL string) error {
	return s.repo.UpdateFotoAnggota(context.Background(), id, fotoURL)
}

func (s *service) GetAnggotaStats(id string) (*AnggotaStats, error) {
	return s.repo.GetAnggotaStats(context.Background(), id)
}

// ─── Pelatih ─────────────────────────────────────────────────────────────────

func (s *service) ListPelatih(cabangID string) ([]Pelatih, error) {
	return s.repo.ListPelatih(context.Background(), cabangID)
}

func (s *service) CreatePelatih(req CreatePelatihRequest) (string, error) {
	return s.repo.CreatePelatih(context.Background(), req)
}

func (s *service) UpdatePelatih(id string, req CreatePelatihRequest) error {
	return s.repo.UpdatePelatih(context.Background(), id, req)
}

func (s *service) DeletePelatih(id string) error {
	return s.repo.DeletePelatih(context.Background(), id)
}

func (s *service) ListPelatihOld(cabangID string) ([]Pelatih, error) {
	return s.repo.ListPelatih(context.Background(), cabangID)
}

func (s *service) GetSebaranProvinsi() ([]SebaranProvinsi, error) {
	return s.repo.GetSebaranProvinsi(context.Background())
}

func (s *service) GetDashboardStats() (*DashboardStats, error) {
	return s.repo.GetDashboardStats(context.Background())
}


func (s *service) UpdateAnggotaKebugaran(id string, req UpdateKebugaranRequest) error {
	return s.repo.UpdateAnggotaKebugaran(context.Background(), id, req)
}

func (s *service) ListPengurus(cabangID string) ([]PengurusCabang, error) {
	return s.repo.ListPengurus(context.Background(), cabangID)
}

func (s *service) CreatePengurus(req CreatePengurusRequest) (string, error) {
	return s.repo.CreatePengurus(context.Background(), req)
}

func (s *service) UpdatePengurus(id string, req CreatePengurusRequest) error {
	return s.repo.UpdatePengurus(context.Background(), id, req)
}

func (s *service) DeletePengurus(id string) error {
	return s.repo.DeletePengurus(context.Background(), id)
}

func (s *service) GetAnggotaKebugaranHistory(id string) ([]map[string]interface{}, error) {
	return s.repo.GetAnggotaKebugaranHistory(context.Background(), id)
}
