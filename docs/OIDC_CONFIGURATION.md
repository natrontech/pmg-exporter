# OIDC Configuration Guide

This guide explains how to configure OIDC authentication with different providers and customize claim names for groups and roles.

## Environment Variables

### Portal Configuration

The portal uses SvelteKit's public environment variables (prefixed with `PUBLIC_`):

```bash
# Required OIDC Configuration
PUBLIC_OIDC_ISSUER=https://your-oidc-provider.com
PUBLIC_OIDC_CLIENT_ID=your-client-id
PUBLIC_BASE_URL=http://localhost:5173

# Optional: Configurable Claims (defaults shown)
PUBLIC_OIDC_GROUPS_CLAIM=groups
PUBLIC_OIDC_ROLES_CLAIM=roles

# Other
PUBLIC_GRAPHQL_ENDPOINT=http://localhost:8080/graphql
PUBLIC_OIDC_PROVIDER_NAME="Your Provider Name"
```

### Orchestrator Configuration

```bash
# Required OIDC Configuration
OIDC_ISSUER=https://your-oidc-provider.com
OIDC_CLIENT_ID=your-client-id

# Optional: Configurable Claims (defaults shown)
OIDC_GROUPS_CLAIM=groups
OIDC_ROLES_CLAIM=roles
OIDC_SCOPE_CLAIM=scope

# Other
GRAPHQL_PORT=8080
ALLOWED_ORIGINS=http://localhost:5173
LOG_LEVEL=info
```

## Provider-Specific Configurations

### Microsoft Azure AD / Entra ID

Azure AD uses different claim names depending on configuration:

```bash
# Option 1: Using groups claim
OIDC_GROUPS_CLAIM=groups
PUBLIC_OIDC_GROUPS_CLAIM=groups

# Option 2: Using roles claim (if configured in app registration)
OIDC_GROUPS_CLAIM=roles
PUBLIC_OIDC_GROUPS_CLAIM=roles
OIDC_ROLES_CLAIM=roles
PUBLIC_OIDC_ROLES_CLAIM=roles
```

**OIDC Issuer Example:**
```bash
OIDC_ISSUER=https://login.microsoftonline.com/{tenant-id}/v2.0
PUBLIC_OIDC_ISSUER=https://login.microsoftonline.com/{tenant-id}/v2.0
```

### Auth0

Auth0 typically uses namespaced claims:

```bash
OIDC_GROUPS_CLAIM=https://yourapp.com/groups
PUBLIC_OIDC_GROUPS_CLAIM=https://yourapp.com/groups
OIDC_ROLES_CLAIM=https://yourapp.com/roles
PUBLIC_OIDC_ROLES_CLAIM=https://yourapp.com/roles
```

**OIDC Issuer Example:**
```bash
OIDC_ISSUER=https://your-domain.auth0.com/
PUBLIC_OIDC_ISSUER=https://your-domain.auth0.com/
```

### Keycloak

Keycloak provides multiple options for groups and roles:

```bash
# Option 1: Groups claim
OIDC_GROUPS_CLAIM=groups
PUBLIC_OIDC_GROUPS_CLAIM=groups

# Option 2: Realm roles
OIDC_ROLES_CLAIM=realm_access.roles
PUBLIC_OIDC_ROLES_CLAIM=realm_access.roles

# Option 3: Client roles (replace 'your-client' with actual client ID)
OIDC_ROLES_CLAIM=resource_access.your-client.roles
PUBLIC_OIDC_ROLES_CLAIM=resource_access.your-client.roles
```

**OIDC Issuer Example:**
```bash
OIDC_ISSUER=https://your-keycloak.com/realms/your-realm
PUBLIC_OIDC_ISSUER=https://your-keycloak.com/realms/your-realm
```

### Okta

Okta uses standard claims:

```bash
OIDC_GROUPS_CLAIM=groups
PUBLIC_OIDC_GROUPS_CLAIM=groups
```

**OIDC Issuer Example:**
```bash
OIDC_ISSUER=https://your-domain.okta.com/oauth2/default
PUBLIC_OIDC_ISSUER=https://your-domain.okta.com/oauth2/default
```

### AWS Cognito

AWS Cognito uses prefixed group claims:

```bash
OIDC_GROUPS_CLAIM=cognito:groups
PUBLIC_OIDC_GROUPS_CLAIM=cognito:groups
```

**OIDC Issuer Example:**
```bash
OIDC_ISSUER=https://cognito-idp.{region}.amazonaws.com/{user-pool-id}
PUBLIC_OIDC_ISSUER=https://cognito-idp.{region}.amazonaws.com/{user-pool-id}
```

### Google

Google Workspace doesn't typically provide groups in standard tokens. You might need custom claims or directory API integration.

```bash
# If using custom claims
OIDC_GROUPS_CLAIM=https://your-domain.com/groups
PUBLIC_OIDC_GROUPS_CLAIM=https://your-domain.com/groups
```

**OIDC Issuer Example:**
```bash
OIDC_ISSUER=https://accounts.google.com
PUBLIC_OIDC_ISSUER=https://accounts.google.com
```

## How It Works

### Token Flow

1. **Authentication**: User logs in via OIDC provider
2. **Token Reception**: Portal receives both `id_token` and `access_token`
3. **Token Storage**: Both tokens are stored as HTTP-only cookies
4. **API Calls**: Portal uses `access_token` for GraphQL API authorization
5. **Authorization**: Orchestrator verifies `access_token` and extracts groups/roles

### Claim Extraction

The system supports flexible claim extraction:

- **Arrays**: `["admin", "user"]` → `["admin", "user"]`
- **Space-separated strings**: `"admin user"` → `["admin", "user"]`
- **Nested objects**: For complex structures like Keycloak's `realm_access.roles`

### Security Considerations

1. **Token Verification**: Both `id_token` and `access_token` are cryptographically verified
2. **Audience Validation**: Tokens are validated against the expected client ID
3. **Issuer Validation**: Tokens are validated against the configured issuer
4. **Signature Verification**: Uses JWKS from the OIDC discovery endpoint

## Troubleshooting

### Common Issues

1. **Missing Groups/Roles**: Check if your OIDC provider includes these claims in the access token
2. **Wrong Claim Name**: Verify the exact claim name used by your provider
3. **Token Audience Mismatch**: Ensure `OIDC_CLIENT_ID` matches your OIDC app registration
4. **CORS Issues**: Ensure `ALLOWED_ORIGINS` includes your portal URL

### Debug Logging

Enable debug logging to see extracted claims:

```bash
LOG_LEVEL=debug
```

This will log the extracted groups and roles for debugging purposes.

## Example Docker Compose

```yaml
version: '3.8'
services:
  orchestrator:
    image: your-orchestrator:latest
    environment:
      OIDC_ISSUER: https://your-provider.com
      OIDC_CLIENT_ID: your-client-id
      OIDC_GROUPS_CLAIM: groups
      OIDC_ROLES_CLAIM: roles
      
  portal:
    image: your-portal:latest
    environment:
      PUBLIC_OIDC_ISSUER: https://your-provider.com
      PUBLIC_OIDC_CLIENT_ID: your-client-id
      PUBLIC_OIDC_GROUPS_CLAIM: groups
      PUBLIC_OIDC_ROLES_CLAIM: roles
      PUBLIC_BASE_URL: https://your-domain.com
      PUBLIC_GRAPHQL_ENDPOINT: https://your-domain.com/graphql
``` 