package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

// Joinked from https://github.com/jacoz/gqlgen-constraint-directive

type constraintDirective struct {
	validator *validator.Validate
}

func New() *constraintDirective {
	return &constraintDirective{
		validator: validator.New(),
	}
}

func (b *constraintDirective) Constraint(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		panic(err)
	}

	path := graphql.GetPathContext(ctx).Path()

	if err = b.validator.Var(val, constraint); err != nil {
		// TODO: Add proper custom graphl error (similar to 'token expired' or 'forbidden').
		return val, fmt.Errorf("value '%s' for %s does not match %s", val, path, constraint)
	}

	return val, nil
}
