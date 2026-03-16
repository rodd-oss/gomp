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

package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/assets"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"image/color"
	"math"
)

type CreateWallManagers struct {
	EntityManager *ecs.EntityManager
	Positions     *stdcomponents.PositionComponentManager
	Rotations     *stdcomponents.RotationComponentManager
	Scales        *stdcomponents.ScaleComponentManager
	BoxColliders  *stdcomponents.BoxColliderComponentManager
	Sprites       *stdcomponents.SpriteComponentManager
	RigidBodies   *stdcomponents.RigidBodyComponentManager
	Renderables   *stdcomponents.RenderableComponentManager
	Velocities    *stdcomponents.VelocityComponentManager
	WallTags      *components.WallTagComponentManager
}

func CreateWall(
	props *CreateWallManagers,
	posX, posY float32,
	angle float64,
	width, height float32,
) ecs.Entity {
	entity := props.EntityManager.Create()
	props.Positions.Create(entity, stdcomponents.Position{
		XY: vectors.Vec2{
			X: posX,
			Y: posY,
		},
	})
	props.Rotations.Create(entity, stdcomponents.Rotation{}.SetFromDegrees(angle))
	props.Scales.Create(entity, stdcomponents.Scale{
		XY: vectors.Vec2{
			X: 1,
			Y: 1,
		},
	})
	props.BoxColliders.Create(entity, stdcomponents.BoxCollider{
		WH: vectors.Vec2{
			X: width,
			Y: height,
		},
		Offset: vectors.Vec2{
			X: 0,
			Y: 0,
		},
		Layer: config.WallCollisionLayer,
		Mask:  0,
	})
	props.RigidBodies.Create(entity, stdcomponents.RigidBody{
		IsStatic: true,
		Mass:     math.MaxFloat32,
	})
	props.Velocities.Create(entity, stdcomponents.Velocity{
		X: 0,
		Y: 0,
	})
	props.Sprites.Create(entity, stdcomponents.Sprite{
		Texture: assets.Textures.Get("wall.png"),
		Frame: rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  width,
			Height: height,
		},
		Origin: rl.Vector2{
			X: 0,
			Y: 0,
		},
		Tint: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
	})
	props.WallTags.Create(entity, components.Wall{})
	props.Renderables.Create(entity, stdcomponents.Renderable{
		Type:       stdcomponents.SpriteRenderableType,
		CameraMask: config.MainCameraLayer | config.MinimapCameraLayer,
	})

	return entity
}
