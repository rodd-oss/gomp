/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package entities

import (
	"gomp_game/cmd/raylib-ecs/assets"
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PlayerStateIdle components.AnimatorState = iota
	PlayerStateWalk
	PlayerStateJump
	PlayerStateFall
	PlayerStateAttack
	PlayerStateHurt
	PlayerStateDie
)

type Player struct {
	ecs.Entity
	Position        *components.Position
	Rotation        *components.Rotation
	Scale           *components.Scale
	SpriteMatrix    *components.SpriteMatrix
	Tint            *components.Tint
	AnimationPlayer *components.AnimationPlayer
}

func CreatePlayer(world *ecs.World) (player Player) {
	// Getting managers
	spriteMatrixes := components.SpriteMatrixService.GetManager(world)
	positions := components.PositionService.GetManager(world)
	rotations := components.RotationService.GetManager(world)
	scales := components.ScaleService.GetManager(world)
	animations := components.AnimationService.GetManager(world)
	tints := components.TintService.GetManager(world)

	// Creating new player

	entity := world.CreateEntity("player")
	player.Entity = entity

	// Adding position component
	t := components.Position{}
	player.Position = positions.Create(entity, t)

	// Adding rotation component
	rotation := components.Rotation{}
	player.Rotation = rotations.Create(entity, rotation)

	// Adding scale component
	scale := components.Scale{
		X: 1,
		Y: 1,
	}
	player.Scale = scales.Create(entity, scale)

	// Adding Tint component
	tint := components.Tint{R: 255, G: 255, B: 255, A: 255}
	player.Tint = tints.Create(entity, tint)

	// Adding sprite matrix component
	spritematrix := components.SpriteMatrix{
		Texture: assets.Textures.Get("assets/milansheet.png"),
		Origin:  rl.Vector2{X: 0.5, Y: 0.5},
		FPS:     12,
		Animations: []components.SpriteMatrixAnimation{
			{
				Name:        "idle",
				Frame:       rl.Rectangle{X: 0, Y: 0, Width: 96, Height: 128},
				NumOfFrames: 1,
				Vertical:    false,
			},
			{
				Name:        "walk",
				Frame:       rl.Rectangle{X: 0, Y: 512, Width: 96, Height: 128},
				NumOfFrames: 8,
				Vertical:    false,
			},
			{
				Name:        "jump",
				Frame:       rl.Rectangle{X: 96, Y: 0, Width: 96, Height: 128},
				NumOfFrames: 1,
				Vertical:    false,
			},
		},
	}
	player.SpriteMatrix = spriteMatrixes.Create(entity, spritematrix)

	// Adding animation component
	animation := components.AnimationPlayer{}
	player.AnimationPlayer = animations.Create(entity, animation)

	return player
}
