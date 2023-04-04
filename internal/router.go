package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WellKnownURLs struct {
	Modules   string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
}

func NewRouter() *chi.Mux {
	m := &moduleHandler{
		store: NewStore(),
	}

	r := chi.NewRouter()
	r.Get("/.well-known/terraform.json", wellKnownTerraformHandler)
	r.Mount("/v1/modules", m.Router())
	return r
}

func wellKnownTerraformHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	addr := WellKnownURLs{
		Modules: fmt.Sprintf("https://%s/v1/modules/", r.Host),
		// Providers: fmt.Sprintf("https://%s/v1/providers/", r.Host),
	}
	if err := json.NewEncoder(w).Encode(addr); err != nil {
		panic(err)
	}
}
