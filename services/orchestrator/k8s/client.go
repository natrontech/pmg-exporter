package k8s

import "k8s.io/client-go/kubernetes"

type Client struct {
	clientset *kubernetes.Clientset
}

func New(clientset *kubernetes.Clientset) (*Client, error) {
	return &Client{clientset: clientset}, nil
}
