package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Handlers holds dependencies for HTTP handlers.
type Handlers struct {
	K8sClient *kubernetes.Clientset
	// Add other dependencies like Helm client, config, etc., if needed by handlers
}

// New creates a new Handlers instance.
func New(k8sClient *kubernetes.Clientset) *Handlers {
	return &Handlers{
		K8sClient: k8sClient,
	}
}

// ListNamespaces retrieves and returns a list of Kubernetes namespaces.
func (h *Handlers) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	// Example: Accessing user info (assuming auth.Context exists in request context)
	/*
	   authCtx, ok := auth.Context(r.Context())
	   if !ok || authCtx.Role < auth.RoleUser { // Or check specific claims/roles
	       http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	       return
	   }
	   slog.Info("ListNamespaces called", "user", authCtx.Claims.Subject)
	*/

	namespaces, err := h.K8sClient.CoreV1().Namespaces().List(r.Context(), metav1.ListOptions{})
	if err != nil {
		slog.Error("failed to list namespaces", "error", err)
		http.Error(w, "Failed to list namespaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(namespaces.Items); err != nil {
		// Log error, but response header is already sent
		slog.Error("failed to encode namespaces response", "error", err)
	}
}

// Add more handler functions as methods on Handlers here...
