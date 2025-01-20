/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package components

import (
	"gomp_game/pkgs/gomp/ecs"
)

const (
	INVALID_ID ecs.ComponentID = iota
	COLOR_ID
	TRANSFORM_ID
	HEALTH_ID
	DESTROY_ID
	SPRITE_ID
	TEXTURE_ID
)
