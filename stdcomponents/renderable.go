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

import "gomp/pkg/ecs"

const (
	InvalidRenderableType RenderableType = iota
	SpriteRenderableType
	SpriteMatrixRenderableType
)

type Renderable struct {
	CameraMask CameraLayer
	Type       RenderableType
	Observed   bool
}

type RenderableType uint8

type RenderableComponentManager = ecs.ComponentManager[Renderable]

func NewRenderableComponentManager() RenderableComponentManager {
	return ecs.NewComponentManager[Renderable](RenderableComponentId)
}
