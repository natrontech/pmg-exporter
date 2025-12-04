package helm

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
)

type Release = release.Release

type Client struct {
	ociRegistry     *url.URL
	credentialsFile string
	kubeconfig      []byte

	config   *action.Configuration
	settings *cli.EnvSettings
}

type ClientOption func(*Client) error

func New(namespace string, options ...ClientOption) (*Client, error) {
	var err error

	client := &Client{}

	for _, opt := range options {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	client.config, client.settings, err = actionConfig(namespace, client.credentialsFile, client.kubeconfig)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func WithCredentialsFile(credentialsFile string) ClientOption {
	return func(c *Client) error {
		c.credentialsFile = credentialsFile
		return nil
	}
}

func WithOCIRegistry(registry string) ClientOption {
	return func(c *Client) error {
		u, err := url.Parse(registry)

		if err != nil {
			return err
		}

		c.ociRegistry = u

		return nil
	}
}

func WithKubeconfig(kubeconfig []byte) ClientOption {
	return func(s *Client) error {
		s.kubeconfig = make([]byte, len(kubeconfig))
		copy(s.kubeconfig, kubeconfig)

		return nil
	}
}

func actionConfig(namespace, credentialsFile string, kubeconfig []byte) (*action.Configuration, *cli.EnvSettings, error) {
	settings := cli.New()
	settings.SetNamespace(namespace)

	// Set up registry client
	opts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptWriter(os.Stdout),
	}

	if credentialsFile != "" {
		opts = append(opts, registry.ClientOptCredentialsFile(credentialsFile))
	}

	registryClient, err := registry.NewClient(opts...)

	if err != nil {
		return nil, nil, err
	}

	getter := settings.RESTClientGetter()

	if len(kubeconfig) > 0 {
		getter = newRESTClientGetter(namespace, kubeconfig)
	}

	logger := func(format string, v ...interface{}) { slog.Debug("helm client", "output", fmt.Sprintf(format, v...)) }
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(getter, namespace, helmDriver, logger); err != nil {
		return nil, nil, err
	}

	actionConfig.RegistryClient = registryClient
	return actionConfig, settings, nil
}
