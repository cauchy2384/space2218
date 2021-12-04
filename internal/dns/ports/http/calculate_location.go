package http

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"space2218/internal/dns"

	"go.uber.org/zap"
)

type LocationCalculator interface {
	CalculateLocation(context.Context, dns.Coordinates3D, dns.Velocity) (dns.Location, error)
}

// swagger:parameters location
type CalculateLocationRequestBody struct {
	// Request
	//
	// in: body
	Request CalculateLocationRequest
}

type CalculateLocationRequest struct {
	// X coordinate
	//
	// required: true
	// example: 123.12
	X string `json:"x"`

	// Y coordinate
	//
	// required: true
	// example: 456.56
	Y string `json:"y"`

	// Z coordinate
	//
	// required: true
	// example: 789.89
	Z string `json:"z"`

	// Velocity
	//
	// required: true
	// example: 20.0
	Velocity string `json:"vel"`
}

func (req CalculateLocationRequest) Parse() (coordinates dns.Coordinates3D, velocity dns.Velocity, err error) {

	coordinates.X, err = dns.CoordinateFromString(req.X)
	if err != nil {
		return dns.Coordinates3D{}, 0, err
	}

	coordinates.Y, err = dns.CoordinateFromString(req.Y)
	if err != nil {
		return dns.Coordinates3D{}, 0, err
	}

	coordinates.Z, err = dns.CoordinateFromString(req.Z)
	if err != nil {
		return dns.Coordinates3D{}, 0, err
	}

	velocity, err = dns.VelocityFromString(req.Velocity)
	if err != nil {
		return dns.Coordinates3D{}, 0, err
	}

	return coordinates, velocity, nil
}

// Location response
// swagger:response locationResponse
type CalculateLocationResponseBody struct {
	// Response
	//
	// in: body
	Response CalculateLocationResponse
}

type CalculateLocationResponse struct {
	// Location
	//
	// example: 1389.57
	Location float64 `json:"loc"`
}

func NewCalculateLocationResponse(l dns.Location) CalculateLocationResponse {
	const roundMultiplier = 100
	return CalculateLocationResponse{
		Location: math.Floor(l.Float64()*roundMultiplier) / roundMultiplier, // round to 2 decimal places
	}
}

func CalculateLocationHandler(logger *zap.Logger, calculator LocationCalculator,
) func(w http.ResponseWriter, r *http.Request) {

	logger = logger.Named("calculate_location")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Info("request received", zap.String("method", r.Method), zap.String("url", r.URL.String()))

		var req CalculateLocationRequest

		if r.Body == nil {
			logger.Info("request body is empty")
			http.Error(w, "empty body", http.StatusUnprocessableEntity)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Info("request body is not valid json", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		coordinates, velocity, err := req.Parse()
		if err != nil {
			logger.Info("request body is not valid", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		location, err := calculator.CalculateLocation(ctx, coordinates, velocity)
		switch {
		case err == nil:
		// ok
		case errors.Is(err, context.Canceled):
			logger.Info("request cancelled", zap.Error(err))
			return
		case errors.Is(err, context.DeadlineExceeded):
			logger.Warn("request timed out", zap.Error(err))
			return
		case errors.Is(err, dns.ErrInvalidValue):
			logger.Info("request body is not valid", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			logger.Error("unexpected error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := NewCalculateLocationResponse(location)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			logger.Error("failed to encode response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("request processed", zap.String("method", r.Method), zap.String("url", r.URL.String()))
	}
}
