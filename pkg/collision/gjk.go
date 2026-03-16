/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

Thank you for your support!
*/

package gjk

import (
	"gomp/stdcomponents"
	"gomp/vectors"
)

const (
	maxIterations = 64
)

func New() GJK {
	return GJK{}
}

type GJK struct {
	simplex Simplex2d
}

type AnyCollider interface {
	GetSupport(direction vectors.Vec2, transform stdcomponents.Transform2d) vectors.Vec2
}

/*
CheckCollision - GJK, Distance, Closest Points
https://www.youtube.com/watch?v=Qupqu1xe7Io
https://dyn4j.org/2010/04/gjk-distance-closest-points/#gjk-distance
*/
func (s *GJK) CheckCollision(
	a, b AnyCollider,
	transformA, transformB stdcomponents.Transform2d,
) bool {
	direction := vectors.Vec2{X: 1, Y: 0}

	p := s.minkowskiSupport2d(a, b, transformA, transformB, direction)
	s.simplex.add(p.ToVec3())
	direction = p.Neg()

	for range maxIterations {
		p = s.minkowskiSupport2d(a, b, transformA, transformB, direction)

		if p.Dot(direction) < 0 {
			return false
		}

		s.simplex.add(p.ToVec3())

		if s.simplex.do(&direction) {
			return true
		}
	}

	panic("GJK infinite loop")
}

func (s *GJK) minkowskiSupport2d(a, b AnyCollider, transformA, transformB stdcomponents.Transform2d, direction vectors.Vec2) vectors.Vec2 {
	aSupport := a.GetSupport(direction, transformA)
	bSupport := b.GetSupport(direction.Neg(), transformB)
	support := aSupport.Sub(bSupport)
	return support
}
