package api

import "context"

type Namespace struct {
	Name string
}

type KubernetesManager interface {
	CreateNamespace(ctx context.Context, name string) (*Namespace, error)
	DeleteNamespace(ctx context.Context, name string) (*Namespace, error)
	Namespaces(ctx context.Context) ([]Namespace, error)
}
