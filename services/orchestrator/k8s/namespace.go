package k8s

import (
	"context"
	api "koda/services/orchestrator"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Client) Namespaces(ctx context.Context) ([]api.Namespace, error) {
	namespaces, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespacesModel := make([]api.Namespace, len(namespaces.Items))
	for i, namespace := range namespaces.Items {
		namespacesModel[i] = api.Namespace{
			Name: namespace.Name,
		}
	}

	return namespacesModel, nil
}

func (c *Client) CreateNamespace(ctx context.Context, name string) (*api.Namespace, error) {
	// TODO: Implement namespace creation logic using c.clientset
	return nil, nil
}

func (c *Client) DeleteNamespace(ctx context.Context, name string) (*api.Namespace, error) {
	// TODO: Implement namespace deletion logic using c.clientset
	return nil, nil
}
