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
	"image/color"
)

type FrameBuffer2D struct {
	Position  rl.Vector2
	Frame     rl.Rectangle
	Texture   rl.RenderTexture2D
	Layer     CameraLayer
	BlendMode rl.BlendMode
	Rotation  float32
	Tint      color.RGBA
	Dst       rl.Rectangle
}

func (d FrameBuffer2D) FrameRect() vectors.Rectangle {
	return vectors.Rectangle{
		X:      d.Frame.X,
		Y:      d.Frame.Y,
		Width:  d.Frame.Width,
		Height: d.Frame.Height,
	}
}

type FrameBuffer2DComponentManager = ecs.ComponentManager[FrameBuffer2D]

func NewFrameBuffer2DComponentManager() FrameBuffer2DComponentManager {
	return ecs.NewComponentManager[FrameBuffer2D](FrameBuffer2DComponentId)
}
