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
	"gomp/vectors"
)

// Camera2D type, defines a 2d camera
type Camera2D struct {
	// Camera offset (displacement from target)
	Offset vectors.Vec2
	// Camera target (rotation and zoom origin)
	Target vectors.Vec2
	// Camera rotation in degrees
	Rotation float32
	// Camera zoom (scaling), should be 1.0f by default
	Zoom float32
}
type Camera struct {
	Camera2D
}

// TODO: remove when raylib will be desintegrated
func (c Camera) ToRaylibCamera() rl.Camera2D {
	return rl.Camera2D{
		Offset:   rl.NewVector2(c.Offset.X, c.Offset.Y),
		Target:   rl.NewVector2(c.Target.X, c.Target.Y),
		Rotation: c.Rotation,
		Zoom:     c.Zoom,
	}
}

type MainCameraComponentManager = ecs.ComponentManager[Camera]

func NewMainCameraComponentManager() MainCameraComponentManager {
	return ecs.NewComponentManager[Camera](MainCameraComponentId)
}

type PipCamera Camera

// TODO: remove when raylib will be desintegrated
func (c PipCamera) ToRaylibCamera() rl.Camera2D {
	return rl.Camera2D{
		Offset:   rl.NewVector2(c.Offset.X, c.Offset.Y),
		Target:   rl.NewVector2(c.Target.X, c.Target.Y),
		Rotation: c.Rotation,
		Zoom:     c.Zoom,
	}
}

type PipCameraComponentManager = ecs.ComponentManager[PipCamera]

func NewPipCameraComponentManager() PipCameraComponentManager {
	return ecs.NewComponentManager[PipCamera](PipCameraComponentId)
}

type MinimapCamera Camera

// TODO: remove when raylib will be desintegrated
func (c MinimapCamera) ToRaylibCamera() rl.Camera2D {
	return rl.Camera2D{
		Offset:   rl.NewVector2(c.Offset.X, c.Offset.Y),
		Target:   rl.NewVector2(c.Target.X, c.Target.Y),
		Rotation: c.Rotation,
		Zoom:     c.Zoom,
	}
}

type MinimapCameraComponentManager = ecs.ComponentManager[MinimapCamera]

func NewMinimapCameraComponentManager() MinimapCameraComponentManager {
	return ecs.NewComponentManager[MinimapCamera](MinimapCameraComponentId)
}
