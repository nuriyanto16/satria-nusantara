package event

// Simplified Event structure for fast implementation
type Event struct {
	ID        string `json:"id"`
	Jenis     string `json:"jenis"` // latgab, ekt, pelatnas
	Nama      string `json:"nama"`
	Tanggal   string `json:"tanggal"`
	Lokasi    string `json:"lokasi"`
	Deskripsi string `json:"deskripsi"`
	Status    string `json:"status"`
}
