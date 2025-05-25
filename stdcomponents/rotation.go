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
	"gomp/vectors"
	"math"
)

//go:generate go tool component -std
type Rotation struct {
	Angle vectors.Radians
}

func (r Rotation) SetFromDegrees(deg float64) Rotation {
	r.Angle = deg * math.Pi / 180
	return r
}

func (r Rotation) Degrees() float64 {
	return r.Angle * 180 / math.Pi
}
