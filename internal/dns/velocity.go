package dns

import (
	"math"
	"strconv"
)

type Velocity float64

func (v Velocity) Float64() float64 {
	return float64(v)
}

func VelocityFromString(s string) (Velocity, error) {
	const bitSize = 64
	f, err := strconv.ParseFloat(s, bitSize)
	return Velocity(f), err
}

func (v Velocity) IsValid() bool {
	return !math.IsInf(v.Float64(), 0) && !math.IsNaN(v.Float64())
}
