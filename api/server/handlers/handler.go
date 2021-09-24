package handlers

import (
	"fmt"
	"net/http"

	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/repository"
)

type PorterHandler interface {
	Config() *config.Config
	Repo() repository.Repository
	HandleAPIError(w http.ResponseWriter, r *http.Request, err apierrors.RequestError)
	HandleAPIErrorNoWrite(w http.ResponseWriter, r *http.Request, err apierrors.RequestError)
	PopulateOAuthSession(w http.ResponseWriter, r *http.Request, state string, isProject bool) error
}

type PorterHandlerWriter interface {
	PorterHandler
	shared.ResultWriter
}

type PorterHandlerReader interface {
	PorterHandler
	shared.RequestDecoderValidator
}

type PorterHandlerReadWriter interface {
	PorterHandlerWriter
	PorterHandlerReader
}

type DefaultPorterHandler struct {
	config           *config.Config
	decoderValidator shared.RequestDecoderValidator
	writer           shared.ResultWriter
}

func NewDefaultPorterHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) PorterHandlerReadWriter {
	return &DefaultPorterHandler{config, decoderValidator, writer}
}

func (d *DefaultPorterHandler) Config() *config.Config {
	return d.config
}

func (d *DefaultPorterHandler) Repo() repository.Repository {
	return d.config.Repo
}

func (d *DefaultPorterHandler) HandleAPIError(w http.ResponseWriter, r *http.Request, err apierrors.RequestError) {
	apierrors.HandleAPIError(d.Config(), w, r, err, true)
}

func (d *DefaultPorterHandler) HandleAPIErrorNoWrite(w http.ResponseWriter, r *http.Request, err apierrors.RequestError) {
	apierrors.HandleAPIError(d.Config(), w, r, err, false)
}

func (d *DefaultPorterHandler) WriteResult(w http.ResponseWriter, r *http.Request, v interface{}) {
	d.writer.WriteResult(w, r, v)
}

func (d *DefaultPorterHandler) DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	return d.decoderValidator.DecodeAndValidate(w, r, v)
}

func (d *DefaultPorterHandler) DecodeAndValidateNoWrite(r *http.Request, v interface{}) error {
	return d.decoderValidator.DecodeAndValidateNoWrite(r, v)
}

func IgnoreAPIError(w http.ResponseWriter, r *http.Request, err apierrors.RequestError) {
	return
}

func (d *DefaultPorterHandler) PopulateOAuthSession(w http.ResponseWriter, r *http.Request, state string, isProject bool) error {
	session, err := d.Config().Store.Get(r, d.Config().ServerConf.CookieName)

	if err != nil {
		return err
	}

	// need state parameter to validate when redirected
	session.Values["state"] = state
	session.Values["query_params"] = r.URL.RawQuery

	if isProject {
		project, _ := r.Context().Value(types.ProjectScope).(*models.Project)

		if project == nil {
			return fmt.Errorf("could not read project")
		}

		session.Values["project_id"] = project.ID
	}

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}