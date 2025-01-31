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

package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"time"
)

// Business

type Health struct {
	Hp, MaxHp int32
}

// Default

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

// Render

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

// Network

type NetworkId int32
type Network struct {
	Id       NetworkId
	PatchIn  []byte
	PatchOut []byte
}
