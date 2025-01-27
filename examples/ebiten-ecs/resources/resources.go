/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package resources

import (
	"embed"
	"gomp/pkgs/gomp"
)

//go:embed **/*.png
var fs embed.FS

var Sprites = gomp.CreateSpriteResource(fs, "sprites")
