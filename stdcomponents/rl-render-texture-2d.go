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
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
)

type RenderTexture2D struct {
	Texture  rl.RenderTexture2D
	Position rl.Vector2
	Frame    rl.Rectangle
}

type RenderTexture2DComponentManager = ecs.ComponentManager[RenderTexture2D]

func NewRenderTexture2DComponentManager() RenderTexture2DComponentManager {
	return ecs.NewComponentManager[RenderTexture2D](RenderTexture2DComponentId)
}
