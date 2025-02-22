package actions

import (
	"loan_service/public"
	"loan_service/templates"

	"github.com/gobuffalo/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		// Use default layout for regular pages
		HTMLLayout: "application.plush.html",

		// But allow overriding it for specific templates
		TemplatesFS: templates.FS(),
		AssetsFS:    public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{},
	})
}
