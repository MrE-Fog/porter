package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	v1 "k8s.io/api/core/v1"
)

// GetK8sNamespacesResponse is the list of namespaces returned when a
// user has successfully authenticated
type GetK8sNamespacesResponse v1.NamespaceList

// GetK8sNamespaces gets a namespaces list in a k8s cluster
func (c *Client) GetK8sNamespaces(
	ctx context.Context,
	projectID uint,
	serviceAccountID uint,
	clusterID uint,
) (*GetK8sNamespacesResponse, error) {
	sa := fmt.Sprintf("%d", serviceAccountID)
	cl := fmt.Sprintf("%d", clusterID)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/projects/%d/k8s/namespaces?"+url.Values{
			"service_account_id": []string{sa},
			"cluster_id":         []string{cl},
		}.Encode(), c.BaseURL, projectID),
		nil,
	)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	bodyResp := &GetK8sNamespacesResponse{}

	if httpErr, err := c.sendRequest(req, bodyResp, true); httpErr != nil || err != nil {
		if httpErr != nil {
			return nil, fmt.Errorf("code %d, errors %v", httpErr.Code, httpErr.Errors)
		}

		return nil, err
	}

	return bodyResp, nil
}
