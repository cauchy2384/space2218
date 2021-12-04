package services

import (
	"context"
	"fmt"
	"space2218/internal/dns"
)

type Navigation struct {
	sectorID dns.SectorID
}

func NewNavigation(sectorID dns.SectorID) *Navigation {
	return &Navigation{
		sectorID: sectorID,
	}
}

func (s *Navigation) CalculateLocation(ctx context.Context,
	coordinates dns.Coordinates3D, velocity dns.Velocity,
) (dns.Location, error) {

	if err := ctx.Err(); err != nil {
		return 0, err
	}

	if !coordinates.IsValid() {
		return 0, fmt.Errorf("coordinates: %w", dns.ErrInvalidValue)
	}
	if !velocity.IsValid() {
		return 0, fmt.Errorf("velocity: %w", dns.ErrInvalidValue)
	}

	location := coordinates.X.Float64() + coordinates.Y.Float64() + coordinates.Z.Float64()
	location = location*s.sectorID.Float64() + velocity.Float64()

	return dns.LocationFromFloat64(location), nil
}
