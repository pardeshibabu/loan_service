package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("home/index.plush.html"))
}

// LoanLifecycleHandler shows the loan lifecycle page
func LoanLifecycleHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("home/loan_lifecycle.plush.html"))
}
