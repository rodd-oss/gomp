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

package vectors

import (
	"math"
)

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

func (v Vec2) Mul(other Vec2) Vec2 {
	return Vec2{v.X * other.X, v.Y * other.Y}
}

func (v Vec2) Div(other Vec2) Vec2 {
	return Vec2{v.X / other.X, v.Y / other.Y}
}

func (v Vec2) AddScalar(scalar float32) Vec2 {
	return Vec2{v.X + scalar, v.Y + scalar}
}

func (v Vec2) SubScalar(scalar float32) Vec2 {
	return Vec2{v.X - scalar, v.Y - scalar}
}

func (v Vec2) Scale(scalar float32) Vec2 {
	return Vec2{v.X * scalar, v.Y * scalar}
}

func (v Vec2) Angle() Radians {
	return math.Atan2(float64(v.Y), float64(v.X))
}

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vec2) Distance(other Vec2) float32 {
	return v.Sub(other).Length()
}

func (v Vec2) Normalize() Vec2 {
	return v.Scale(1 / v.Length())
}

func (v Vec2) Rotate(angle Radians) Vec2 {
	return Vec2{
		v.X*float32(math.Cos(angle)) - v.Y*float32(math.Sin(angle)),
		v.X*float32(math.Sin(angle)) + v.Y*float32(math.Cos(angle)),
	}
}

func (v Vec2) LengthSquared() float32 {
	l := v.Length()
	return l * l
}

func (v Vec2) Neg() Vec2 {
	return Vec2{-v.X, -v.Y}
}

// Perpendicular - clockwise
func (v Vec2) Perpendicular() Vec2 {
	return Vec2{
		X: v.Y,
		Y: -v.X,
	}
}

func (v Vec2) Normal() Vec2 {
	return Vec2{
		X: v.Y,
		Y: -v.X,
	}
}

func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) ToVec3() Vec3 {
	return Vec3{v.X, v.Y, 0}
}
