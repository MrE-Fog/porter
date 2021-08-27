package authz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/porter-dev/porter/api/server/authz/policy"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
)

type PolicyMiddleware struct {
	config       *config.Config
	endpointMeta types.APIRequestMetadata
	loader       policy.PolicyDocumentLoader
}

func NewPolicyMiddleware(
	config *config.Config,
	endpointMeta types.APIRequestMetadata,
	loader policy.PolicyDocumentLoader,
) *PolicyMiddleware {
	return &PolicyMiddleware{config, endpointMeta, loader}
}

func (p *PolicyMiddleware) Middleware(next http.Handler) http.Handler {
	return &PolicyHandler{next, p.config, p.endpointMeta, p.loader}
}

type PolicyHandler struct {
	next         http.Handler
	config       *config.Config
	endpointMeta types.APIRequestMetadata
	loader       policy.PolicyDocumentLoader
}

func (h *PolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the full map of scopes to resource actions
	reqScopes, reqErr := getRequestActionForEndpoint(r, h.endpointMeta)

	if reqErr != nil {
		apierrors.HandleAPIError(r.Context(), h.config, w, reqErr)
		return
	}

	// load policy documents for the user + project
	projID := reqScopes[types.ProjectScope].Resource.UInt
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	policyDocs, reqErr := h.loader.LoadPolicyDocuments(user.ID, projID)

	if reqErr != nil {
		apierrors.HandleAPIError(r.Context(), h.config, w, reqErr)
		return
	}

	// validate that the policy permits the action
	hasAccess := policy.HasScopeAccess(policyDocs, reqScopes)

	if !hasAccess {
		apierrors.HandleAPIError(
			r.Context(),
			h.config,
			w,
			apierrors.NewErrForbidden(fmt.Errorf("policy forbids action for user %d in project %d", user.ID, projID)),
		)

		return
	}

	// add the set of resource ids to the request context
	ctx := NewRequestScopeCtx(r.Context(), reqScopes)
	r = r.WithContext(ctx)
	h.next.ServeHTTP(w, r)
}

const RequestScopeCtxKey = "requestscopes"

func NewRequestScopeCtx(ctx context.Context, reqScopes map[types.PermissionScope]*policy.RequestAction) context.Context {
	return context.WithValue(ctx, RequestScopeCtxKey, reqScopes)
}

func getRequestActionForEndpoint(
	r *http.Request,
	endpointMeta types.APIRequestMetadata,
) (res map[types.PermissionScope]*policy.RequestAction, reqErr apierrors.RequestError) {
	res = make(map[types.PermissionScope]*policy.RequestAction)

	// iterate through scopes, attach policies as needed
	for _, scope := range endpointMeta.Scopes {
		// find the resource ID and create the resource
		resource := types.NameOrUInt{}

		switch scope {
		case types.ProjectScope:
			resource.UInt, reqErr = GetURLParamUint(r, string(types.URLParamProjectID))
		case types.ClusterScope:
			resource.UInt, reqErr = GetURLParamUint(r, string(types.URLParamClusterID))
		case types.RegistryScope:
			resource.UInt, reqErr = GetURLParamUint(r, string(types.URLParamRegistryID))
		case types.HelmRepoScope:
			resource.UInt, reqErr = GetURLParamUint(r, string(types.URLParamHelmRepoID))
		case types.GitInstallationScope:
			resource.UInt, reqErr = GetURLParamUint(r, string(types.URLParamGitInstallationID))
		case types.InfraScope:
			resource.UInt, reqErr = GetURLParamUint(r, string(types.URLParamInfraID))
		case types.NamespaceScope:
			resource.Name, reqErr = GetURLParamString(r, string(types.URLParamNamespace))
		case types.ReleaseScope:
			resource.Name, reqErr = GetURLParamString(r, string(types.URLParamReleaseName))
		}

		if reqErr != nil {
			return nil, reqErr
		}

		res[scope] = &policy.RequestAction{
			Verb:     endpointMeta.Verb,
			Resource: resource,
		}
	}

	return res, nil
}
