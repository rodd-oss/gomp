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

type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v Vec3) Mul(other Vec3) Vec3 {
	return Vec3{v.X * other.X, v.Y * other.Y, v.Z * other.Z}
}

func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{
		a.Y*b.Z - b.Y*a.Z,
		a.Z*b.X - b.Z*a.X,
		a.X*b.Y - b.X*a.Y,
	}
}

func (v Vec3) Div(other Vec3) Vec3 {
	return Vec3{v.X / other.X, v.Y / other.Y, v.Z / other.Z}
}

func (v Vec3) AddScalar(scalar float32) Vec3 {
	return Vec3{v.X + scalar, v.Y + scalar, v.Z + scalar}
}

func (v Vec3) SubScalar(scalar float32) Vec3 {
	return Vec3{v.X - scalar, v.Y - scalar, v.Z - scalar}
}

func (v Vec3) Scale(scalar float32) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vec3) Normalize() Vec3 {
	return v.Scale(1 / v.Length())
}

func (v Vec3) Rotate(angle Radians) Vec3 {
	return Vec3{
		v.X*float32(math.Cos(angle)) - v.Y*float32(math.Sin(angle)),
		v.X*float32(math.Sin(angle)) + v.Y*float32(math.Cos(angle)),
		v.Z,
	}
}

func (v Vec3) LengthSquared() float32 {
	l := v.Length()
	return l * l
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Perpendicular() Vec3 {
	return Vec3{
		X: -v.Y,
		Y: v.X,
	}
}

func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vec3) ToVec2() Vec2 {
	return Vec2{v.X, v.Y}
}
