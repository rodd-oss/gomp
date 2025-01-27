/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package components

import (
	"gomp/pkgs/ecs"
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position struct {
	X, Y float32
}
type Rotation struct {
	Angle float32
}
type Scale struct {
	X, Y float32
}
type Mirrored struct {
	X, Y bool
}
type Health struct {
	Hp, MaxHp int32
}
type Sprite struct {
	Texture *rl.Texture2D
	Frame   rl.Rectangle
	Origin  rl.Vector2
	Tint    color.RGBA
}
type SpriteSheet struct {
	Texture     *rl.Texture2D
	Frame       rl.Rectangle
	Origin      rl.Vector2
	NumOfFrames int32
	FPS         int32
	Vertical    bool
}
type Tint = color.RGBA
type SpriteMatrixAnimation struct {
	Name        string
	Frame       rl.Rectangle
	NumOfFrames uint8
	Vertical    bool
	Loop        bool
}
type SpriteMatrix struct {
	Texture    *rl.Texture2D
	Origin     rl.Vector2
	FPS        int32
	Animations []SpriteMatrixAnimation
}
type TextureRender struct {
	Texture  *rl.Texture2D
	Frame    rl.Rectangle
	Origin   rl.Vector2
	Tint     color.RGBA
	Dest     rl.Rectangle
	Rotation float32
}
type AnimationPlayer struct {
	First         uint8
	Last          uint8
	Current       uint8
	Speed         float32
	Loop          bool
	Vertical      bool
	ElapsedTime   time.Duration
	FrameDuration time.Duration
	State         AnimationState
	IsInitialized bool
}

type AnimationState int

var PositionService = ecs.CreateComponentService[Position](POSITION_ID)
var RotationService = ecs.CreateComponentService[Rotation](ROTATION_ID)
var ScaleService = ecs.CreateComponentService[Scale](SCALE_ID)
var MirroredService = ecs.CreateComponentService[Mirrored](MIRRORED_ID)

var HealthService = ecs.CreateComponentService[Health](HEALTH_ID)

var SpriteService = ecs.CreateComponentService[Sprite](SPRITE_ID)
var SpriteSheetService = ecs.CreateComponentService[SpriteSheet](SPRITE_SHEET_ID)
var SpriteMatrixService = ecs.CreateComponentService[SpriteMatrix](SPRITE_MATRIX_ID)
var TintService = ecs.CreateComponentService[Tint](TINT_ID)

var AnimationPlayerService = ecs.CreateComponentService[AnimationPlayer](ANIMATION_ID)
var AnimationStateService = ecs.CreateComponentService[AnimationState](ANIMATION_STATE_ID)

var TextureRenderService = ecs.CreateComponentService[TextureRender](TEXTURE_RENDER_ID)

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has Color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// Destroy system
