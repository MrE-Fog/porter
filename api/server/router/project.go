package router

import (
	"github.com/go-chi/chi"
	"github.com/porter-dev/porter/api/server/handlers/cluster"
	"github.com/porter-dev/porter/api/server/handlers/project"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/types"
)

func NewProjectScopedRegisterer(children ...*Registerer) *Registerer {
	return &Registerer{
		GetRoutes: GetProjectScopedRoutes,
		Children:  children,
	}
}

func GetProjectScopedRoutes(
	r chi.Router,
	config *shared.Config,
	basePath *types.Path,
	factory shared.APIEndpointFactory,
	children ...*Registerer,
) []*Route {
	routes, projPath := getProjectRoutes(r, config, basePath, factory)

	if len(children) > 0 {
		r.Route(projPath.RelativePath, func(r chi.Router) {
			for _, child := range children {
				childRoutes := child.GetRoutes(r, config, basePath, factory, child.Children...)

				routes = append(routes, childRoutes...)
			}
		})
	}

	return routes
}

func getProjectRoutes(
	r chi.Router,
	config *shared.Config,
	basePath *types.Path,
	factory shared.APIEndpointFactory,
) ([]*Route, *types.Path) {
	relPath := "/projects/{project_id}"

	newPath := &types.Path{
		Parent:       basePath,
		RelativePath: relPath,
	}

	routes := make([]*Route, 0)

	// GET /api/projects/{project_id} -> project.NewProjectGetHandler
	getEndpoint := factory.NewAPIEndpoint(
		&types.APIRequestMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &types.Path{
				Parent:       basePath,
				RelativePath: relPath,
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.ProjectScope,
			},
		},
	)

	getHandler := project.NewProjectGetHandler(
		config,
		factory.GetResultWriter(),
	)

	routes = append(routes, &Route{
		Endpoint: getEndpoint,
		Handler:  getHandler,
		Router:   r,
	})

	// GET /api/projects/{project_id}/clusters -> cluster.NewClusterListHandler
	listClusterEndpoint := factory.NewAPIEndpoint(
		&types.APIRequestMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &types.Path{
				Parent:       basePath,
				RelativePath: relPath + "/clusters",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.ProjectScope,
			},
		},
	)

	listClusterHandler := cluster.NewClusterListHandler(
		config,
		factory.GetResultWriter(),
	)

	routes = append(routes, &Route{
		Endpoint: listClusterEndpoint,
		Handler:  listClusterHandler,
		Router:   r,
	})

	return routes, newPath
}
