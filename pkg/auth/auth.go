package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type TokenClaims struct {
	Subject  string                 `json:"sub,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Username string                 `json:"preferred_username,omitempty"`
	Email    string                 `json:"email,omitempty"`
	Groups   []string               `json:"groups,omitempty"`
	Roles    []string               `json:"roles,omitempty"`
	Scope    string                 `json:"scope,omitempty"`
	Scopes   []string               `json:"scopes,omitempty"`
	Claims   map[string]interface{} `json:"-"`
	jwt.RegisteredClaims
}

type ClaimsConfig struct {
	GroupsClaim string
	RolesClaim  string
	ScopeClaim  string
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func Middleware(clientID, issuer string, claimsConfig ClaimsConfig) (mux.MiddlewareFunc, error) {
	if issuer == "" || clientID == "" {
		return nil, errors.New("clientID and issuer are required to initialize auth middleware")
	}

	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authContext := authContext(r.Header.Get("Authorization"), clientID, issuer, provider, claimsConfig)
			ctx := setAuthContext(r.Context(), &authContext)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

func authContext(authHeader, clientID, issuer string, provider *oidc.Provider, claimsConfig ClaimsConfig) AuthContext {
	if authHeader == "" {
		return AuthContext{
			Role: RoleAnonymous,
		}
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		slog.Error("invalid authorization header", "header", authHeader)
		return AuthContext{
			Role: RoleAnonymous,
		}
	}

	accessToken := parts[1]

	claims, err := verifyAccessToken(accessToken, issuer, claimsConfig)
	if err != nil {
		slog.Error("access token verification failed", "error", err)
		return AuthContext{
			Role: RoleAnonymous,
		}
	}

	if claims.Subject == "" {
		slog.Error("sub claim is not defined in access token")
		return AuthContext{
			Role: RoleAnonymous,
		}
	}

	return AuthContext{
		Role:   RoleUser,
		Claims: claims,
	}
}

func verifyAccessToken(tokenString, issuer string, claimsConfig ClaimsConfig) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found in token header")
		}

		jwks, err := fetchJWKS(issuer)
		if err != nil {
			return nil, err
		}

		return getPublicKey(kid, jwks)
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	claims := &TokenClaims{}

	if sub, ok := mapClaims["sub"].(string); ok {
		claims.Subject = sub
	}
	if name, ok := mapClaims["name"].(string); ok {
		claims.Name = name
	}
	if username, ok := mapClaims["preferred_username"].(string); ok {
		claims.Username = username
	}
	if email, ok := mapClaims["email"].(string); ok {
		claims.Email = email
	}

	extractCustomClaims(claims, mapClaims, claimsConfig)

	if claims.Scopes == nil && claims.Scope != "" {
		claims.Scopes = strings.Fields(claims.Scope)
	}

	return claims, nil
}

func extractCustomClaims(claims *TokenClaims, rawClaims jwt.MapClaims, claimsConfig ClaimsConfig) {
	if groupsClaim, exists := rawClaims[claimsConfig.GroupsClaim]; exists {
		if groupsArray, ok := groupsClaim.([]interface{}); ok {
			groups := make([]string, 0, len(groupsArray))
			for _, g := range groupsArray {
				if groupStr, ok := g.(string); ok {
					groups = append(groups, groupStr)
				}
			}
			claims.Groups = groups
		} else if groupsStr, ok := groupsClaim.(string); ok {
			claims.Groups = strings.Fields(groupsStr)
		}
	}

	if rolesClaim, exists := rawClaims[claimsConfig.RolesClaim]; exists {
		if rolesArray, ok := rolesClaim.([]interface{}); ok {
			roles := make([]string, 0, len(rolesArray))
			for _, r := range rolesArray {
				if roleStr, ok := r.(string); ok {
					roles = append(roles, roleStr)
				}
			}
			claims.Roles = roles
		} else if rolesStr, ok := rolesClaim.(string); ok {
			claims.Roles = strings.Fields(rolesStr)
		}
	}

	if scopeClaim, exists := rawClaims[claimsConfig.ScopeClaim]; exists {
		if scopeStr, ok := scopeClaim.(string); ok {
			claims.Scope = scopeStr
		}
	}
}

func fetchJWKS(issuer string) (*JWKSet, error) {
	resp, err := http.Get(fmt.Sprintf("%s/.well-known/openid-configuration", strings.TrimSuffix(issuer, "/")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	jwksURI, ok := config["jwks_uri"].(string)
	if !ok {
		return nil, errors.New("jwks_uri not found in discovery document")
	}

	resp, err = http.Get(jwksURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func getPublicKey(kid string, jwks *JWKSet) (*rsa.PublicKey, error) {
	var jwk *JWK
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			jwk = &key
			break
		}
	}

	if jwk == nil {
		return nil, errors.New("key not found")
	}

	if jwk.Kty != "RSA" {
		return nil, errors.New("only RSA keys are supported")
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := big.NewInt(0).SetBytes(nBytes)
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}
