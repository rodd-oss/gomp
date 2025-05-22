/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package stdcomponents

import (
	"image/color"
)

//go:generate go tool component -std
type Tint struct {
	R, G, B, A uint8
}

func (t Tint) RGBA() color.RGBA {
	return color.RGBA{t.R, t.G, t.B, t.A}
}
