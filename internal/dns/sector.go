package dns

type SectorID uint64

func (s SectorID) Float64() float64 {
	return float64(s)
}
