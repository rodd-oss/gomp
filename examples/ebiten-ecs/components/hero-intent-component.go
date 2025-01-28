/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package components

import (
	"gomp/pkg/gomp"

	"github.com/yohamta/donburi/features/math"
)

type HeroIntentData struct {
	Move math.Vec2
	Jump bool
	Fire bool
}

var HeroIntentComponent = gomp.CreateComponent[HeroIntentData]()
