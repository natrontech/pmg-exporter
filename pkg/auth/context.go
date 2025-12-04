package auth

import "context"

type contextKey string

const claimsKey contextKey = "authClaims"
const roleKey contextKey = "authRole"
const authContextKey contextKey = "authContext"

type AuthContext struct {
	Claims    *TokenClaims
	Role      UserRole
	IsExpired bool
}

func Claims(ctx context.Context) (*TokenClaims, bool) {
	claims := ctx.Value(claimsKey)

	if claims == nil {
		return nil, false
	}

	result, ok := claims.(*TokenClaims)

	if !ok {
		return nil, false
	}

	return result, true
}

func Role(ctx context.Context) (UserRole, bool) {
	role := ctx.Value(roleKey)

	if role == nil {
		return RoleAnonymous, false
	}

	result, ok := role.(UserRole)

	if !ok {
		return RoleAnonymous, false
	}

	return result, true
}

func Context(ctx context.Context) (*AuthContext, bool) {
	authContext := ctx.Value(authContextKey)

	if authContext == nil {
		return nil, false
	}

	result, ok := authContext.(*AuthContext)

	if !ok {
		return nil, false
	}

	return result, true
}

func setAuthContext(ctx context.Context, auth *AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, auth)
}
