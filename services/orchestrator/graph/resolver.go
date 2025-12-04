package graph

import (
	api "koda/services/orchestrator"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	Kubernetes api.KubernetesManager
}
