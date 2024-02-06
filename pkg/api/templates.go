package api

// Template represents a template
// TODO: this implementation expects a small amount of templates. TODO: we should use a nosql document store.
type Template struct {
	ID        string `json:"template_id"`
	Template  string `json:"template"`
	IsEnabled bool   `json:"is_enabled"`
}
