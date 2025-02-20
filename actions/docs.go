package actions

import (
	"net/http"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
)

// DocsHandler serves the Swagger UI
func DocsHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("docs/swagger.html"))
}

// DocsFS returns the docs filesystem
func DocsFS() http.FileSystem {
	return http.Dir(filepath.Join("docs"))
}
