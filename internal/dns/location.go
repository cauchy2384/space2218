package dns

type Location float64

func LocationFromFloat64(v float64) Location {
	return Location(v)
}

func (l Location) Float64() float64 {
	return float64(l)
}
