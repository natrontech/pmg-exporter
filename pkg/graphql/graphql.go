package graphql

import (
	"context"
	"log/slog"

	"koda/pkg/auth"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

var publicFields = []string{
	"__schema",
	"__type",
	"_service",
	"_entities",
	"__typename",
}

func Authorize(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	operationContext := graphql.GetOperationContext(ctx)
	selectionSet := operationContext.Operation.SelectionSet

	authContext, ok := auth.Context(ctx)

	if !ok {
		slog.Info("user is not authenticated")
		return unauthenticatedResponse
	}

	if authContext.IsExpired {
		slog.Info("expired token")
		return expiredTokenResponse
	}

	if !isAuthorized(operationContext, selectionSet, authContext.Role, nil) {
		slog.Info("user is not authorized for this operation", "role", authContext.Role)
		return forbiddenResponse
	}

	return next(ctx)
}

func isAuthorized(ctx *graphql.OperationContext, selection ast.SelectionSet, userRole auth.UserRole, requiredParentRoles []auth.UserRole) bool {
	for _, selection := range graphql.CollectFields(ctx, selection, nil) {
		if selection.Field == nil {
			return false
		}

		// Allow introspection query without authentication.
		if isPublic(selection.Field) {
			continue
		}

		requiredRoles := requiredRoles(selection.Field)

		if len(requiredRoles) == 0 {
			if len(requiredParentRoles) == 0 {
				slog.Info("field has no hasRole directive and no parent with hasRole directive", "field", selection.Field.Name)
				return false
			}

			// Children inherit roles from parent if not explicitly specified.
			requiredRoles = requiredParentRoles
		}

		if !isFieldAuthorized(selection.Field, userRole, requiredRoles) {
			return false
		}

		return isAuthorized(ctx, selection.SelectionSet, userRole, requiredRoles)
	}

	return true
}

func isPublic(field *ast.Field) bool {
	for _, f := range publicFields {
		if field.Name == f {
			return true
		}
	}

	return false
}

func isFieldAuthorized(field *ast.Field, userRole auth.UserRole, requiredRoles []auth.UserRole) bool {
	if len(requiredRoles) == 0 {
		slog.Info("hasRole directive has no valid role specified", "field", field.Name)
		return false
	}

	for _, r := range requiredRoles {
		if userRole.IsAuthorized(r) {
			// If user is authorized for one of the required roles, access is granted.
			return true
		}
	}

	return false
}

func requiredRoles(field *ast.Field) []auth.UserRole {
	hasRoleDir := field.Definition.Directives.ForName("hasRole")

	if hasRoleDir == nil {
		return nil
	}

	roleArgs := hasRoleDir.Arguments.ForName("role")

	if roleArgs == nil {
		return nil
	}

	roles := make([]auth.UserRole, 0)

	for _, roleArg := range roleArgs.Value.Children {
		if roleArg.Value.Raw == "" {
			continue
		}

		r, err := auth.NewRole(roleArg.Value.Raw)

		if err != nil {
			slog.Error("parsing argument from GraphQL hasRole directive", "error", err)
			continue
		}

		roles = append(roles, r)
	}

	return roles
}
