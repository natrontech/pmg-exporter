package router

import (
	"net/http"

	"koda/services/orchestrator/handlers" // Import the handlers package

	"github.com/gorilla/mux"
)

// New creates and configures a mux router.
func New(
	h *handlers.Handlers, // Pass Handlers struct
	healthHandler http.HandlerFunc, // Pass the specific health handler func
	authMiddleware mux.MiddlewareFunc, // Pass the auth middleware
) *mux.Router {
	r := mux.NewRouter()

	// IMPORTANT: Register unprotected routes BEFORE the auth middleware
	r.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)

	// Apply auth middleware globally. All routes registered AFTER this are protected.
	r.Use(authMiddleware)

	// API routes (protected by global authMiddleware)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/namespaces", h.ListNamespaces).Methods(http.MethodGet)
	// Add more API routes here, calling methods on h...
	// e.g., api.HandleFunc("/pods", h.ListPods).Methods(http.MethodGet)

	return r
}
