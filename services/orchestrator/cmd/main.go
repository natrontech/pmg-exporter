package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"koda/pkg/auth"
	"koda/pkg/graceful"
	"koda/pkg/logger"
	"koda/services/orchestrator/config"
	"koda/services/orchestrator/graph"
	"koda/services/orchestrator/health"
	"koda/services/orchestrator/k8s"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var cfg config.Config

	if err := envconfig.Process("", &cfg); err != nil {
		slog.Error("Failed to process environment variables", "error", err)
		os.Exit(1)
	}

	logger.SetLevel(cfg.LogLevel)

	ctx, cancel := graceful.Context()
	defer cancel()

	claimsConfig := auth.ClaimsConfig{
		GroupsClaim: cfg.OIDCGroupsClaim,
		RolesClaim:  cfg.OIDCRolesClaim,
		ScopeClaim:  cfg.OIDCScopeClaim,
	}

	authMiddleware, err := auth.Middleware(cfg.OIDCClientID, cfg.OIDCIssuer, claimsConfig)
	if err != nil {
		slog.Error("Failed to create auth middleware", "error", err)
		os.Exit(1)
	}

	// ociRegistry := fmt.Sprintf("oci://%s", cfg.Registry)
	// helmOpts := []helm.ClientOption{helm.WithOCIRegistry(ociRegistry)}

	// if cfg.RegistryCredentialsFile != "" {
	// 	helmOpts = append(helmOpts, helm.WithCredentialsFile(cfg.RegistryCredentialsFile))
	// }

	// helmClient, err := helm.New("default", helmOpts...)

	if err != nil {
		slog.Error("failed to initialize helm client", "error", err)
		os.Exit(1)
	}

	var k8sClient *kubernetes.Clientset

	k8sClient, err = newKubernetesClient()

	if err != nil {
		slog.Error("failed to initialize kubernetes client", "error", err)
		os.Exit(1)
	}

	healthChecker := health.New()
	for name, check := range health.DefaultChecks() {
		healthChecker.AddCheck(name, check)
	}

	healthChecker.AddCheck("kubernetes", health.KubernetesCheck(k8sClient))

	kubernetesManager, err := k8s.New(
		k8sClient,
	)

	if err != nil {
		slog.Error("failed to initialize kubernetes client", "error", err)
		os.Exit(1)
	}

	resolver := graph.Resolver{
		Kubernetes: kubernetesManager,
	}

	graphqlServer, err := NewGraphQLServer(cfg.GraphQLPort, cfg.AllowedOrigins, &resolver, []mux.MiddlewareFunc{authMiddleware}, healthChecker)

	if err != nil {
		slog.Error("failed to initialize graphql server", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting GraphQL server", "port", cfg.GraphQLPort, "version", cfg.Version)
	graceful.Serve(ctx, graphqlServer)
}

func newKubernetesClient() (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config
	var kubeconfig string

	if envVar := os.Getenv("KUBECONFIG"); envVar != "" {
		kubeconfig = envVar
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		config, err = rest.InClusterConfig()

		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return clientset, nil
}
