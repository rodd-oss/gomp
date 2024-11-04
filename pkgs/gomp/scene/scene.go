/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package scene

import "gomp_game/pkgs/gomp/ecs"

type Scene struct {
	Name string

	Systems  []ecs.System
	Entities []ecs.Entity
}
