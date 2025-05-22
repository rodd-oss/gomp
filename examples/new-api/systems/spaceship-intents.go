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

package systems

import (
	"gomp/examples/new-api/assets"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/entities"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
	"time"
)

func NewSpaceshipIntentsSystem() SpaceshipIntentsSystem {
	return SpaceshipIntentsSystem{}
}

type SpaceshipIntentsSystem struct {
	EntityManager         *ecs.EntityManager
	SpaceshipIntents      *components.SpaceshipIntentComponentManager
	Positions             *stdcomponents.PositionComponentManager
	Velocities            *stdcomponents.VelocityComponentManager
	Rotations             *stdcomponents.RotationComponentManager
	Scales                *stdcomponents.ScaleComponentManager
	BoxColliders          *stdcomponents.BoxColliderComponentManager
	CircleColliders       *stdcomponents.CircleColliderComponentManager
	BulletTags            *components.BulletTagComponentManager
	Sprites               *stdcomponents.SpriteComponentManager
	Textures              *stdcomponents.RLTextureProComponentManager
	RigidBodies           *stdcomponents.RigidBodyComponentManager
	Weapons               *components.WeaponComponentManager
	Hps                   *components.HpComponentManager
	SoundEffects          *components.SoundEffectComponentManager
	TexturePositionSmooth *stdcomponents.TexturePositionSmoothComponentManager
	Renderables           *stdcomponents.RenderableComponentManager
	RenderOrders          *stdcomponents.RenderOrderComponentManager
	moveSpeed             float32
}

func (s *SpaceshipIntentsSystem) Init() {}
func (s *SpaceshipIntentsSystem) Run(dt time.Duration) {
	var moveSpeedMax float32 = 300
	var moveSpeedMaxBackwards float32 = -200
	var rotateSpeed vectors.Radians = 3
	var speedIncrement float32 = 10

	var bulletSpeed float32 = 300

	dtSec := float32(dt.Seconds())

	s.SpaceshipIntents.EachEntity(func(entity ecs.Entity) bool {
		intent := s.SpaceshipIntents.GetUnsafe(entity)
		vel := s.Velocities.GetUnsafe(entity)
		rot := s.Rotations.GetUnsafe(entity)
		pos := s.Positions.GetUnsafe(entity)
		weapon := s.Weapons.GetUnsafe(entity)
		hp := s.Hps.GetUnsafe(entity)
		flySfx := s.SoundEffects.GetUnsafe(entity)

		if intent.RotateLeft {
			rot.Angle -= rotateSpeed * vectors.Radians(dtSec)
		}
		if intent.RotateRight {
			rot.Angle += rotateSpeed * vectors.Radians(dtSec)
		}
		if intent.MoveUp {
			s.moveSpeed += speedIncrement
			if s.moveSpeed > moveSpeedMax {
				s.moveSpeed = moveSpeedMax
			}
		}
		if intent.MoveDown {
			s.moveSpeed -= speedIncrement
			if s.moveSpeed < moveSpeedMaxBackwards {
				s.moveSpeed = moveSpeedMaxBackwards
			}
		}

		if !intent.MoveUp && !intent.MoveDown {
			if s.moveSpeed > 0 {
				s.moveSpeed -= speedIncrement
			} else if s.moveSpeed < 0 {
				s.moveSpeed += speedIncrement
			}
		}

		absMoveSpeed := math.Abs(float64(s.moveSpeed))
		flySfx.Volume = float32(absMoveSpeed / float64(moveSpeedMax))

		if (intent.MoveUp || intent.MoveDown || intent.RotateLeft || intent.RotateRight || s.moveSpeed != 0) && !flySfx.IsPlaying {
			flySfx.IsPlaying = true
		} else if !(intent.MoveUp || intent.MoveDown || intent.RotateLeft || intent.RotateRight || s.moveSpeed != 0) && flySfx.IsPlaying || hp.Hp == 0 {
			flySfx.IsPlaying = false
		}

		vel.Y = float32(math.Cos(rot.Angle+math.Pi)) * s.moveSpeed
		vel.X = -float32(math.Sin(rot.Angle+math.Pi)) * s.moveSpeed

		if weapon.CooldownLeft <= 0 {
			if intent.Fire {
				var count int = 30
				for i := range count {
					var angle = math.Pi*2/float64(count)*float64(i) + rot.Angle

					bulletVelocityY := vel.Y + float32(math.Cos(angle+math.Pi))*bulletSpeed
					bulletVelocityX := vel.X - float32(math.Sin(angle+math.Pi))*bulletSpeed
					entities.CreateBullet(entities.CreateBulletManagers{
						EntityManager:   s.EntityManager,
						Positions:       s.Positions,
						Rotations:       s.Rotations,
						Scales:          s.Scales,
						Velocities:      s.Velocities,
						CircleColliders: s.CircleColliders,
						BoxColliders:    s.BoxColliders,
						RigidBodies:     s.RigidBodies,
						Sprites:         s.Sprites,
						BulletTags:      s.BulletTags,
						Hps:             s.Hps,
						Renderables:     s.Renderables,
					}, pos.XY.X, pos.XY.Y, angle, bulletVelocityX, bulletVelocityY)
				}
				weapon.CooldownLeft = weapon.Cooldown

				fireSoundEntity := s.EntityManager.Create()

				s.SoundEffects.Create(fireSoundEntity, components.SoundEffect{
					Clip:      assets.Audio.Get("gun_sound.wav"),
					IsPlaying: false,
					IsLooping: false,
					Pitch:     1.0,
					Volume:    1.0,
					Pan:       0.5,
				})

				s.Positions.Create(
					fireSoundEntity,
					stdcomponents.Position{
						XY: vectors.Vec2{
							X: pos.XY.X,
							Y: pos.XY.Y,
						},
					},
				)
			}
		} else {
			weapon.CooldownLeft -= dt
		}

		return true
	})
}

func (s *SpaceshipIntentsSystem) Destroy() {}
