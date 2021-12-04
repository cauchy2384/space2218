package services

import (
	"context"
	"errors"
	"space2218/internal/dns"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNavigation(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		SectorID         dns.SectorID
		Coordinates      dns.Coordinates3D
		Velocity         dns.Velocity
		ExpectedLocation dns.Location
	}{
		{
			SectorID:         1,
			Coordinates:      dns.Coordinates3D{X: 123.12, Y: 456.56, Z: 789.89},
			Velocity:         20.0,
			ExpectedLocation: 1389.57,
		},
	}

	for _, tt := range testCases {
		service := NewNavigation(tt.SectorID)

		location, err := service.CalculateLocation(ctx, tt.Coordinates, tt.Velocity)
		require.NoError(t, err)
		assert.InDelta(t, tt.ExpectedLocation.Float64(), location.Float64(), 0.01)
	}
}

func TestNavigationCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	service := NewNavigation(0)
	_, err := service.CalculateLocation(ctx, dns.Coordinates3D{}, 0)
	assert.True(t, errors.Is(err, context.Canceled))

}
