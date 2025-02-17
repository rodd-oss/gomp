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

package components

import (
	"gomp"
	"gomp/pkg/ecs"
)

type GameComponents struct {
	*gomp.DesktopComponents
	Health *HealthComponentManager
}

func NewGameComponents(world *ecs.World, desktopComponents *gomp.DesktopComponents) GameComponents {
	return GameComponents{
		DesktopComponents: desktopComponents,
		Health:            NewHealthComponentManager(world),
	}
}
