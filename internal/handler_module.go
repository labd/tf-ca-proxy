package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type moduleRequestKey struct{}

type moduleHandler struct {
	store *Store
}

func (m *moduleHandler) Router() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(moduleCtx)
		r.Group(func(r chi.Router) {
			r.Use(authTokenHandler)
			r.Get("/{namespace}/{name}/{provider}/versions", m.versionListHandler)
			r.Get("/{namespace}/{name}/{provider}/{version}", m.versionHandler)
			r.Get("/{namespace}/{name}/{provider}/{version}/download", m.downloadHandler)
		})
		r.Get("/{namespace}/{name}/{provider}/{version}/assets/{asset}", m.assetHandler)
	})
	return r
}

func (m *moduleHandler) versionListHandler(w http.ResponseWriter, r *http.Request) {
	moduleRequest := GetModuleRequest(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	sourceName := fmt.Sprintf("%s/terraform-%s-%s",
		chi.URLParam(r, "namespace"),
		chi.URLParam(r, "provider"),
		chi.URLParam(r, "name"),
	)

	versions, err := m.store.listModuleVersions(r.Context(), moduleRequest)
	if err != nil {
		panic(err)
	}

	result := ModuleVersionResponse{
		Modules: []ModuleData{
			{
				Source:   sourceName,
				Versions: versions,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func (m *moduleHandler) versionHandler(w http.ResponseWriter, r *http.Request) {
	moduleRequest := GetModuleRequest(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	info, err := m.store.getModuleVersion(r.Context(), moduleRequest)
	if err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		panic(err)
	}
}

func (m *moduleHandler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	moduleRequest := GetModuleRequest(r.Context())

	filename, err := m.store.getModuleVersionAssets(r.Context(), moduleRequest)
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("https://%s/v1/modules/%s/%s/%s/%s/assets/%s",
		r.Host,
		moduleRequest.Namespace,
		moduleRequest.Name,
		moduleRequest.Provider,
		moduleRequest.Version,
		filename)

	w.Header().Add("X-Terraform-Get", url)
	w.WriteHeader(http.StatusOK)
}

func (m *moduleHandler) assetHandler(w http.ResponseWriter, r *http.Request) {
	moduleRequest := GetModuleRequest(r.Context())
	asset := chi.URLParam(r, "asset")

	reader, err := m.store.downloadModuleVersion(r.Context(), moduleRequest, asset)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, reader); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WithModuleRequest(r *http.Request) context.Context {
	return context.WithValue(r.Context(), moduleRequestKey{}, ModuleRequest{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
		Version:   chi.URLParam(r, "version"),
	})
}

func GetModuleRequest(ctx context.Context) ModuleRequest {
	return ctx.Value(moduleRequestKey{}).(ModuleRequest)
}

func moduleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithModuleRequest(r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
