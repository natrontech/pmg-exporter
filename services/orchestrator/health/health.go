package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the general health check configuration
type Config struct {
	Timeout time.Duration `envconfig:"HEALTH_TIMEOUT" default:"5s"` // Default timeout for health checks
}

var healthConfig Config

// Status represents the health status of a component
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// Component represents a health check component
type Component struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Check represents a health check function
type Check func(ctx context.Context) Component

// Health manages health checks for the application
type Health struct {
	mu      sync.RWMutex
	checks  map[string]Check
	timeout time.Duration
}

// New creates a new Health instance
func New() *Health {
	_ = envconfig.Process("", &healthConfig) // Process env vars
	return &Health{
		checks:  make(map[string]Check),
		timeout: healthConfig.Timeout,
	}
}

// AddCheck adds a new health check
func (h *Health) AddCheck(name string, check Check) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// RunChecks runs all registered health checks
func (h *Health) RunChecks(ctx context.Context) []Component {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	var components []Component
	for _, check := range h.checks {
		components = append(components, check(ctx))
	}
	return components
}

// IsHealthy returns true if all components are up
func (h *Health) IsHealthy(components []Component) bool {
	for _, component := range components {
		if component.Status == StatusDown {
			return false
		}
	}
	return true
}

// LivenessHandler returns an http.HandlerFunc for liveness probe
func (h *Health) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Liveness probe should be lightweight
		w.WriteHeader(http.StatusOK)
	}
}

// ReadinessHandler returns an http.HandlerFunc for readiness probe
func (h *Health) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		components := h.RunChecks(r.Context())
		if !h.IsHealthy(components) {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// HealthHandler returns an http.HandlerFunc for detailed health check
func (h *Health) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		components := h.RunChecks(r.Context())
		w.Header().Set("Content-Type", "application/json")
		if !h.IsHealthy(components) {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "UP",
			"components": components,
		})
	}
}
