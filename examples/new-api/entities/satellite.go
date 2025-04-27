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
	"gomp/examples/new-api/assets"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/config"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CreateSatelliteManagers struct {
	EntityManager *ecs.EntityManager
	Positions     *stdcomponents.PositionComponentManager
	Rotations     *stdcomponents.RotationComponentManager
	Scales        *stdcomponents.ScaleComponentManager
	Sprites       *stdcomponents.SpriteComponentManager
	BoxColliders  *stdcomponents.BoxColliderComponentManager
	RigidBodies   *stdcomponents.RigidBodyComponentManager
	Velocities    *stdcomponents.VelocityComponentManager
	Renderables   *stdcomponents.RenderableComponentManager
	SoundEffects  *components.SoundEffectsComponentManager
}

func CreateSatellite(
	props CreateSatelliteManagers,
	posX, posY float32,
	angle float64,
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

	props.Velocities.Create(entity, stdcomponents.Velocity{
		X: 0,
		Y: 0,
	})

	props.Sprites.Create(entity, stdcomponents.Sprite{
		Texture: assets.Textures.Get("satellite_B.png"),
		Origin:  rl.Vector2{X: 32, Y: 40},
		Frame:   rl.Rectangle{0, 0, 64, 64},
		Tint:    color.RGBA{255, 255, 255, 255},
	})

	props.BoxColliders.Create(entity, stdcomponents.BoxCollider{
		WH: vectors.Vec2{
			X: 64,
			Y: 64,
		},
		Offset: vectors.Vec2{
			X: 32,
			Y: 32,
		},
		Layer: config.PlayerCollisionLayer,
		Mask:  1<<config.EnemyCollisionLayer | 1<<config.BulletCollisionLayer | 1<<config.PlayerCollisionLayer,
	})
	props.RigidBodies.Create(entity, stdcomponents.RigidBody{
		IsStatic: false,
		Mass:     1,
	})
	props.Renderables.Create(entity, stdcomponents.Renderable{
		Type:       stdcomponents.SpriteRenderableType,
		CameraMask: config.MainCameraLayer | config.MinimapCameraLayer,
	})

	props.SoundEffects.Create(entity, components.SoundEffect{
		Clip:      assets.Audio.Get("satellite_main_sound.wav"),
		IsLooping: true,
		IsPlaying: false,
		Volume:    0.05,
		Pan:       0.5,
		Pitch:     1.0,
	})

	return entity
}
