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
	"math"
)

const (
	Culling2DFullscreenBB uint8 = iota
	Culling2DNone               // TODO: set default?
	// Culling2DDistance
	// Culling2DRectangle
	// Culling2DOcclusion difficult to implement
	// Culling2DFrustum huh
)

type CameraLayer uint64

// Camera component, defines a camera component
// Camera2D is a struct that defines a 2D camera
// Dst defines the camera's destination on the renderer screen
// Layer defines the camera's layer, default is 0 - disabled
// Order defines the camera's order, ascending order
type Camera struct {
	rl.Camera2D
	Dst       vectors.Rectangle // TODO: remove?
	Layer     CameraLayer
	Culling   uint8
	Order     int
	BlendMode rl.BlendMode
	BGColor   color.RGBA
	Tint      color.RGBA
}

func (c *Camera) Rect() vectors.Rectangle {
	// Calculate the non-rotated top-left corner of the view rectangle
	x := c.Target.X - (c.Offset.X / c.Zoom)
	y := c.Target.Y - (c.Offset.Y / c.Zoom)
	width := c.Offset.X * 2 / c.Zoom
	height := c.Offset.Y * 2 / c.Zoom

	// When rotation is zero, we can return directly.
	if c.Rotation == 0 {
		return vectors.Rectangle{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}
	}

	// Define the four corners of the non-rotated rectangle
	topLeft := vectors.Vec2{X: x, Y: y}
	topRight := vectors.Vec2{X: x + width, Y: y}
	bottomRight := vectors.Vec2{X: x + width, Y: y + height}
	bottomLeft := vectors.Vec2{X: x, Y: y + height}

	// Rotate each corner around the camera.Target using the camera rotation
	topLeft = rotatePoint(topLeft, vectors.Vec2(c.Target), float64(c.Rotation))
	topRight = rotatePoint(topRight, vectors.Vec2(c.Target), float64(c.Rotation))
	bottomRight = rotatePoint(bottomRight, vectors.Vec2(c.Target), float64(c.Rotation))
	bottomLeft = rotatePoint(bottomLeft, vectors.Vec2(c.Target), float64(c.Rotation))

	// Determine the axis-aligned bounding box that contains all rotated points
	// TODO: fast math 32bit
	minX := math.Min(math.Min(float64(topLeft.X), float64(topRight.X)), math.Min(float64(bottomRight.X), float64(bottomLeft.X)))
	maxX := math.Max(math.Max(float64(topLeft.X), float64(topRight.X)), math.Max(float64(bottomRight.X), float64(bottomLeft.X)))
	minY := math.Min(math.Min(float64(topLeft.Y), float64(topRight.Y)), math.Min(float64(bottomRight.Y), float64(bottomLeft.Y)))
	maxY := math.Max(math.Max(float64(topLeft.Y), float64(topRight.Y)), math.Max(float64(bottomRight.Y), float64(bottomLeft.Y)))

	return vectors.Rectangle{
		X:      float32(minX),
		Y:      float32(minY),
		Width:  float32(maxX - minX),
		Height: float32(maxY - minY),
	}
}

// rotatePoint rotates point p around pivot by angle degrees.
func rotatePoint(p, pivot vectors.Vec2, angle float64) vectors.Vec2 {
	// Convert angle from degrees to radians.
	// TODO: fast math 32bit
	theta := angle * (math.Pi / 180)
	sinTheta := float32(math.Sin(theta))
	cosTheta := float32(math.Cos(theta))

	// Translate point to origin
	dx := p.X - pivot.X
	dy := p.Y - pivot.Y

	// Apply rotation matrix
	rotatedX := dx*cosTheta - dy*sinTheta
	rotatedY := dx*sinTheta + dy*cosTheta

	// Translate point back
	return vectors.Vec2{
		X: rotatedX + pivot.X,
		Y: rotatedY + pivot.Y,
	}
}

type CameraComponentManager = ecs.ComponentManager[Camera]

func NewCameraComponentManager() CameraComponentManager {
	return ecs.NewComponentManager[Camera](CameraComponentId)
}
