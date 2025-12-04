package health

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesCheck checks if the Kubernetes client is healthy
func KubernetesCheck(k8s *kubernetes.Clientset) Check {
	return func(ctx context.Context) Component {
		if k8s == nil {
			return Component{
				Name:   "kubernetes",
				Status: StatusDown,
				Error:  "Kubernetes client not initialized",
			}
		}

		// Try to list namespaces as a basic connectivity test
		_, err := k8s.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return Component{
				Name:   "kubernetes",
				Status: StatusDown,
				Error:  err.Error(),
			}
		}

		return Component{
			Name:   "kubernetes",
			Status: StatusUp,
		}
	}
}
