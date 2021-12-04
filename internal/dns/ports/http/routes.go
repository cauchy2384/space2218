package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/heptiolabs/healthcheck"
	"go.uber.org/zap"
)

func Routes(logger *zap.Logger, locationCalculator LocationCalculator) *chi.Mux {
	// handlers
	r := chi.NewRouter()

	// healthcheck & readiness
	health := healthcheck.NewHandler()
	healtcheckLogger := logger.Named("healthcheck")
	health.AddLivenessCheck("alive", func() error {
		healtcheckLogger.Info("alive")
		return nil
	})
	health.AddReadinessCheck("ready", func() error {
		healtcheckLogger.Info("ready")
		return nil
	})

	// swagger:route GET /live healthcheck liveness
	//
	// Liveness probe
	//
	//     Responses:
	//       200:
	r.Get("/live", health.LiveEndpoint)

	// swagger:route GET /ready healthcheck readiness
	//
	// Readiness probe
	//
	//     Responses:
	//       200:
	r.Get("/ready", health.ReadyEndpoint)

	// api
	r.Route("/api/v1", func(r chi.Router) {
		// swagger:route POST /location api location
		//
		// Calculate location by given coordinates and velocity
		//
		//	 Responses:
		//	   200: locationResponse
		r.Post("/location", CalculateLocationHandler(logger, locationCalculator))
	})

	return r
}
