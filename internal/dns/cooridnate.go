package dns

import (
	"math"
	"strconv"
)

type Coordinate float64

func (c Coordinate) Float64() float64 {
	return float64(c)
}

func CoordinateFromString(s string) (Coordinate, error) {
	const bitSize = 64
	f, err := strconv.ParseFloat(s, bitSize)
	return Coordinate(f), err
}

func (c Coordinate) IsValid() bool {
	return !math.IsInf(c.Float64(), 0) && !math.IsNaN(c.Float64())
}

type Coordinates3D struct {
	X Coordinate
	Y Coordinate
	Z Coordinate
}

func (c Coordinates3D) IsValid() bool {
	return c.X.IsValid() && c.Y.IsValid() && c.Z.IsValid()
}
