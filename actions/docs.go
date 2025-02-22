package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// DocsHandler serves the API documentation
func DocsHandler(c buffalo.Context) error {
	// Set layout to empty string to prevent using application layout
	c.Set("layout", "")
	return c.Render(http.StatusOK, r.HTML("docs/layout.plush.html"))
}
