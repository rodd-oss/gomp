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

package gjk

import (
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
)

const (
	epaMaxTolerance = 0.01 // should not be very low because of circle colliders
)

/*
EPA - Expanding Polytope Algorithm
Based on https://dyn4j.org/2010/05/epa-expanding-polytope-algorithm/#epa-alternatives
*/
func (s *GJK) EPA(
	a, b AnyCollider,
	transformA, transformB stdcomponents.Transform2d,
) (vectors.Vec2, float32) {
	polytope := s.simplex.toPolytope(make([]vectors.Vec2, 0, 6))

	bestNormal := vectors.Vec2{}
	bestDistance := float32(math.MaxFloat32)
	bestTolerance := float32(math.MaxFloat32)

	for range maxIterations {
		edge := s.findClosestEdge(polytope)
		point := s.minkowskiSupport2d(a, b, transformA, transformB, edge.normal)
		distance := point.Dot(edge.normal)
		tolerance := distance - edge.distance
		if tolerance < epaMaxTolerance {
			return edge.normal, distance
		}
		if tolerance < bestTolerance {
			bestTolerance = tolerance
			bestNormal = edge.normal
			bestDistance = distance
		}

		polytope = append(polytope[:edge.index], append([]vectors.Vec2{point}, polytope[edge.index:]...)...)
	}

	return bestNormal, bestDistance
}

func (s *GJK) findClosestEdge(polytope []vectors.Vec2) closestEdge {
	closest := closestEdge{
		distance: float32(math.MaxFloat32),
		normal:   vectors.Vec2{},
		index:    -1,
	}

	for i := 0; i < len(polytope); i++ {
		j := (i + 1) % len(polytope)
		a := polytope[i]
		b := polytope[j]

		edge := b.Sub(a)
		normal := edge.Perpendicular().Normalize()
		distance := normal.Dot(a)

		if distance < closest.distance {
			closest.distance = distance
			closest.normal = normal
			closest.index = j
		}
	}

	return closest
}

type closestEdge struct {
	distance float32
	normal   vectors.Vec2
	index    int
}
