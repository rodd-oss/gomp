/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdcomponents

import (
	"gomp/pkg/ecs"
)

const (
	InvalidComponentId ecs.ComponentId = 1<<16 - 1 - iota
	PositionComponentId
	RotationComponentId
	ScaleComponentId
	FlipComponentId
	VelocityComponentId
	SpriteComponentId
	SpriteSheetComponentId
	SpriteMatrixComponentId
	TextureRenderComponentId
	AnimationPlayerComponentId
	AnimationStateComponentId
	TintComponentId
	NetworkComponentId
)
