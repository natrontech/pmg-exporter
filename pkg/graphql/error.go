package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrorCode string

const (
	UnauthorizedErrorCode ErrorCode = "UNAUTHORIZED"
	ForbiddenErrorCode    ErrorCode = "FORBIDDEN"
	TokenExpiredErrorCode ErrorCode = "TOKEN_EXPIRED"
)

func unauthenticatedResponse(ctx context.Context) *graphql.Response {
	return errorResponse(UnauthorizedErrorCode, "Unauthorized")
}

func forbiddenResponse(ctx context.Context) *graphql.Response {
	return errorResponse(ForbiddenErrorCode, "Forbidden")
}

func expiredTokenResponse(ctx context.Context) *graphql.Response {
	return errorResponse(TokenExpiredErrorCode, "Token Expired")
}

func errorResponse(code ErrorCode, message string) *graphql.Response {
	return &graphql.Response{
		Errors: []*gqlerror.Error{{
			Message: message,
			Extensions: map[string]interface{}{
				"code": code,
			},
		}},
		Data: nil,
	}
}
