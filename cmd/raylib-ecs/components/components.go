/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package components

import (
	"gomp_game/pkgs/gomp/ecs"
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
	Tint        color.RGBA
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
	First         int32
	Last          int32
	Current       int32
	Speed         float32
	Loop          bool
	Vertical      bool
	ElapsedTime   time.Duration
	FrameDuration time.Duration
}
type AnimatorState int
type Animator struct {
	LastState AnimatorState
	State     AnimatorState
}

var PositionService = ecs.CreateComponentService[Position](POSITION_ID)
var RotationService = ecs.CreateComponentService[Rotation](ROTATION_ID)
var ScaleService = ecs.CreateComponentService[Scale](SCALE_ID)

var HealthService = ecs.CreateComponentService[Health](HEALTH_ID)

var SpriteService = ecs.CreateComponentService[Sprite](SPRITE_ID)
var SpriteSheetService = ecs.CreateComponentService[SpriteSheet](SPRITE_SHEET_ID)
var SpriteMatrixService = ecs.CreateComponentService[SpriteMatrix](SPRITE_MATRIX_ID)
var AnimationService = ecs.CreateComponentService[AnimationPlayer](ANIMATION_ID)
var AnimatorService = ecs.CreateComponentService[Animator](ANIMATOR_ID)
var TintService = ecs.CreateComponentService[Tint](TINT_ID)

var TextureRenderService = ecs.CreateComponentService[TextureRender](TEXTURE_RENDER_ID)

// spawn creature every tick with random hp and position
// each creature looses hp every tick
// each creature has Color that depends on its own maxHP and current hp
// when hp == 0 creature dies

// spawn system
// movement system
// hp system
// Destroy system
