package training

import (
	"context"

	"satria-nusantara/backend/internal/organization"
)

type Service interface {
	ListSesi(unitID string, tanggal string) ([]Sesi, error)
	CreateSesi(req CreateSesiRequest) (string, error)
	UpdateSesi(id string, req CreateSesiRequest) error
	DeleteSesi(id string) error
	GenerateQR(sesiID string, durasiJam int) (*QRAbsensi, error)
	ScanQR(req ScanQRRequest) error
	GetRingkasan(sesiID string) (*RingkasanKehadiran, []Kehadiran, error)
}

type service struct {
	repo    Repository
	orgRepo organization.Repository
}

func NewService(repo Repository, orgRepo organization.Repository) Service {
	return &service{repo: repo, orgRepo: orgRepo}
}

func (s *service) ListSesi(unitID, tanggal string) ([]Sesi, error) {
	return s.repo.ListSesi(context.Background(), unitID, tanggal)
}

func (s *service) CreateSesi(req CreateSesiRequest) (string, error) {
	return s.repo.CreateSesi(context.Background(), req)
}

func (s *service) UpdateSesi(id string, req CreateSesiRequest) error {
	return s.repo.UpdateSesi(context.Background(), id, req)
}

func (s *service) DeleteSesi(id string) error {
	return s.repo.DeleteSesi(context.Background(), id)
}

func (s *service) GenerateQR(sesiID string, durasiJam int) (*QRAbsensi, error) {
	return s.repo.GenerateQR(context.Background(), sesiID, durasiJam)
}

func (s *service) ScanQR(req ScanQRRequest) error {
	return s.repo.ScanQR(context.Background(), req)
}

func (s *service) GetRingkasan(sesiID string) (*RingkasanKehadiran, []Kehadiran, error) {
	ctx := context.Background()
	kehadiran, err := s.repo.GetKehadiran(ctx, sesiID)
	if err != nil {
		return nil, nil, err
	}
	
	// Untuk menghitung persentase, kita butuh tahu total anggota unit tersebut
	// Sederhananya, ini bisa didapat dari sesi -> unit -> count(anggota)
	// Kita bisa asumsikan Sesi punya UnitID, mari hitung tanpa total spesifik (0)
	// Dalam realitasnya kita akan panggil orgRepo.GetUnitByID
	
	ringkasan, err := s.repo.RingkasanKehadiran(ctx, sesiID, 0)
	if err != nil {
		return nil, nil, err
	}
	return ringkasan, kehadiran, nil
}
