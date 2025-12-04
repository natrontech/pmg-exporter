package config

// Config holds application configuration variables.
type Config struct {
	GraphQLPort    string `envconfig:"GRAPHQL_PORT" default:"8080"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:5173"`
	LogLevel       string `envconfig:"LOG_LEVEL" default:"info"`
	OIDCIssuer     string `envconfig:"OIDC_ISSUER"`
	OIDCClientID   string `envconfig:"OIDC_CLIENT_ID"`
	Version        string `envconfig:"VERSION" default:"dev"`

	// OIDC Claims Configuration
	OIDCGroupsClaim string `envconfig:"OIDC_GROUPS_CLAIM" default:"groups"`
	OIDCRolesClaim  string `envconfig:"OIDC_ROLES_CLAIM" default:"roles"`
	OIDCScopeClaim  string `envconfig:"OIDC_SCOPE_CLAIM" default:"scope"`

	RegistryPullSecret      string `envconfig:"CONTAINER_REGISTRY_PULL_SECRET"`
	RegistryCredentialsFile string `envconfig:"CONTAINER_REGISTRY_CONFIG"`
	Registry                string `envconfig:"CONTAINER_REGISTRY" default:"ghcr.io/natrontech/koda"`
}
