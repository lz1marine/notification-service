package api

// Template represents a template
// TODO: this implementation expects a small amount of templates. TODO: we should use a nosql document store.
type Template struct {
	// ID is the id of the template
	ID string `json:"template_id"`

	// Template is the template content
	Template string `json:"template"`

	// IsEnabled indicates if the template is enabled
	IsEnabled bool `json:"is_enabled"`
}
