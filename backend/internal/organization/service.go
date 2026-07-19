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
	VerifikasiAnggota(id string, req VerifikasiAnggotaRequest) error
	UpdateFotoAnggota(id, fotoURL string) error
	GetAnggotaStats(id string) (*AnggotaStats, error)
	UpdateAnggotaKebugaran(id string, req UpdateKebugaranRequest) error

	// Pelatih
	ListPelatih(cabangID string) ([]Pelatih, error)

	// Sebaran
	GetSebaranProvinsi() ([]SebaranProvinsi, error)

	// Dashboard
	GetDashboardStats() (*DashboardStats, error)
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

func (s *service) VerifikasiAnggota(id string, req VerifikasiAnggotaRequest) error {
	if req.Aksi != "approve" && req.Aksi != "reject" {
		return errors.New("aksi harus 'approve' atau 'reject'")
	}
	status := "aktif"
	if req.Aksi == "reject" {
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

func (s *service) GetSebaranProvinsi() ([]SebaranProvinsi, error) {
	return s.repo.GetSebaranProvinsi(context.Background())
}

func (s *service) GetDashboardStats() (*DashboardStats, error) {
	return s.repo.GetDashboardStats(context.Background())
}


func (s *service) UpdateAnggotaKebugaran(id string, req UpdateKebugaranRequest) error {
	return s.repo.UpdateAnggotaKebugaran(context.Background(), id, req)
}
