//go:build !dev

package web

import (
	"embed"
	"text/template"
)

//go:embed src/**/*.html
var HtmlContent embed.FS

func init() {
	UiTemplates = &Template{
		templates: template.Must(template.ParseFS(HtmlContent, "src/**/*.html")),
	}
}
