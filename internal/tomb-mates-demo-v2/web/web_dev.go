//go:build dev

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package web

import (
	"text/template"
)

func init() {
	UiTemplates = &Template{
		templates: template.Must(template.ParseGlob("web/src/**/*.html")),
	}
}
