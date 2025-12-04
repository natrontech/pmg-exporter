package graph

import (
	api "koda/services/orchestrator"
	"koda/services/orchestrator/graph/model"
)

func convertNamespace(namespace api.Namespace) model.Namespace {
	return model.Namespace{
		Name: namespace.Name,
	}
}
