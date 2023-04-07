package internal

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type ModuleRequest struct {
	Namespace string
	Name      string
	Provider  string
	Version   string
}

func (r ModuleRequest) PackageName() string {
	return "terraform-" + r.Provider + "-" + r.Name
}

type ModuleVersion struct {
	Version string `json:"version"`
}

type ModuleData struct {
	Source   string          `json:"source"`
	Versions []ModuleVersion `json:"versions"`
}

type ModuleVersionResponse struct {
	Modules []ModuleData `json:"modules"`
}

type ModuleVersionInfo struct {
	ID          string     `json:"id"`
	Owner       string     `json:"owner"`
	Namespace   string     `json:"namespace"`
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Provider    string     `json:"provider"`
	Tag         string     `json:"tag"`
	PublishedAt *time.Time `json:"published_at"`
}
