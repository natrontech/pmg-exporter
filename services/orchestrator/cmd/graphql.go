package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"koda/pkg/graphql"
	"koda/services/orchestrator/graph"
	"koda/services/orchestrator/graph/directive"
	"koda/services/orchestrator/graph/model"
	"koda/services/orchestrator/health"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type GraphQLServer struct {
	server *http.Server
	port   string
}

func NewGraphQLServer(port string, allowedOriginsConfig string, resolver *graph.Resolver, middlewares []mux.MiddlewareFunc, healthChecker *health.Health) (*GraphQLServer, error) {
	server := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			Constraint: directive.New().Constraint,
			HasRole: func(ctx context.Context, obj interface{}, next gqlgen.Resolver, role []model.Role) (interface{}, error) {
				// No-op directive, since auth logic is enforced in operation middleware graphql.Authorize
				return next(ctx)
			},
		},
	}))

	server.AddTransport(transport.Options{})
	server.AddTransport(transport.POST{})

	server.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	server.Use(extension.Introspection{})

	server.AroundOperations(graphql.Authorize)
	server.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := gqlgen.DefaultErrorPresenter(ctx, e)

		slog.Error("interal server error", "error", err)

		// TODO: Return a generic error like "errors.New("internal server error")" here to not leak internal information
		return err
	})

	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler())
	router.Use(middlewares...)

	router.HandleFunc("/health", healthChecker.HealthHandler())
	router.HandleFunc("/health/live", healthChecker.LivenessHandler())
	router.HandleFunc("/health/ready", healthChecker.ReadinessHandler())

	queryEndpoint := "/graphql"
	router.Handle("/playground", playground.Handler("GraphQL playground", queryEndpoint))
	router.Handle(queryEndpoint, server)

	allowedOrigins := handlers.AllowedOrigins([]string{allowedOriginsConfig})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	allowCredentials := handlers.AllowCredentials()

	corsHandler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders, allowCredentials)(router)

	return &GraphQLServer{
		port: port,
		server: &http.Server{
			Addr:    ":" + port,
			Handler: corsHandler,
		},
	}, nil
}

func (s *GraphQLServer) Start() error {
	slog.Info("graphql playground enabled", "url", fmt.Sprintf("http://localhost:%s/playground", s.port))
	return s.server.ListenAndServe()
}

func (s *GraphQLServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *GraphQLServer) Label() string {
	return "GraphQL"
}
