package finance

// Simplified Finance structure for fast implementation
type BLBA struct {
	ID        string  `json:"id"`
	AnggotaID string  `json:"anggota_id"`
	Periode   string  `json:"periode"`
	Nominal   float64 `json:"nominal"`
	Status    string  `json:"status"` // lunas, tunggakan
}
