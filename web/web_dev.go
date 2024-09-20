//go:build dev

package web

import (
	"text/template"
)

func init() {
	UiTemplates = &Template{
		templates: template.Must(template.ParseGlob("web/src/**/*.html")),
	}
}
