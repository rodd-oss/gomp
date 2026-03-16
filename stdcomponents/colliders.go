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
	"gomp/pkg/ecs"
	"gomp/vectors"
	"math"
)

type ColliderShape uint8

const (
	InvalidColliderShape ColliderShape = iota
	BoxColliderShape
	CircleColliderShape
	PolygonColliderShape
)

type CollisionMask uint64

func (m CollisionMask) HasLayer(layer CollisionLayer) bool {
	return m&(1<<layer) != 0
}

type CollisionLayer = CollisionMask

const (
	CollisionLayerNone CollisionLayer = 0
)

// ===========================
// Box Collider
// ===========================

type BoxCollider struct {
	WH         vectors.Vec2
	Offset     vectors.Vec2
	Layer      CollisionLayer
	Mask       CollisionMask
	AllowSleep bool
}

func (c *BoxCollider) GetSupport(direction vectors.Vec2, transform Transform2d) vectors.Vec2 {
	// Precompute rotation terms once
	cos := float32(math.Cos(transform.Rotation))
	sin := float32(math.Sin(transform.Rotation))

	// Inverse-rotate direction to local space (avoids per-vertex rotation)
	localDir := vectors.Vec2{
		X: direction.X*cos + direction.Y*sin,
		Y: -direction.X*sin + direction.Y*cos,
	}

	// Branchless selection using sign bits (Go-optimized)
	xSign := math.Float32bits(localDir.X) >> 31
	ySign := math.Float32bits(localDir.Y) >> 31
	localSupport := vectors.Vec2{
		X: c.WH.X * (1 - float32(xSign)), // 0 if negative, WH.X otherwise
		Y: c.WH.Y * (1 - float32(ySign)),
	}

	// Apply offset and rotate to world space
	vertex := localSupport.Sub(c.Offset)
	rotated := vectors.Vec2{
		X: vertex.X*cos - vertex.Y*sin,
		Y: vertex.X*sin + vertex.Y*cos,
	}

	return rotated.Mul(transform.Scale).Add(transform.Position)
}

type BoxColliderComponentManager = ecs.ComponentManager[BoxCollider]

func NewBoxColliderComponentManager() BoxColliderComponentManager {
	return ecs.NewComponentManager[BoxCollider](ColliderBoxComponentId)
}

// ===========================
// Circle Collider
// ===========================

type CircleCollider struct {
	Radius     float32
	Layer      CollisionLayer
	Mask       CollisionMask
	Offset     vectors.Vec2
	AllowSleep bool
}

func (c *CircleCollider) GetSupport(direction vectors.Vec2, transform Transform2d) vectors.Vec2 {
	var radiusWithOffset vectors.Vec2
	// Handle zero direction to avoid division by zero
	if direction.X == 0 && direction.Y == 0 {
		// Fallback to a default direction (e.g., right)
		defaultDir := vectors.Vec2{X: c.Radius, Y: 0}
		radiusWithOffset = defaultDir.Sub(c.Offset).Mul(transform.Scale)
	} else {
		// Compute scaled direction without trigonometry
		mag := float32(math.Hypot(float64(direction.X), float64(direction.Y)))
		invMag := c.Radius / mag
		scaledDir := vectors.Vec2{
			X: direction.X * invMag,
			Y: direction.Y * invMag,
		}
		// Apply offset, scale, and translation
		radiusWithOffset = scaledDir.Sub(c.Offset).Mul(transform.Scale)
	}

	return transform.Position.Add(radiusWithOffset)
}

type CircleColliderComponentManager = ecs.ComponentManager[CircleCollider]

func NewCircleColliderComponentManager() CircleColliderComponentManager {
	return ecs.NewComponentManager[CircleCollider](ColliderCircleComponentId)
}

// ===========================
// Polygon Collider
// ===========================

type PolygonCollider struct {
	Vertices   []vectors.Vec2
	Layer      CollisionLayer
	Mask       CollisionMask
	Offset     vectors.Vec2
	AllowSleep bool
}

func (c *PolygonCollider) GetSupport(direction vectors.Vec2, transform Transform2d) vectors.Vec2 {
	maxDot := math.Inf(-1)
	var maxVertex vectors.Vec2

	for _, v := range c.Vertices {
		scaled := v.Mul(transform.Scale)
		rotated := scaled.Rotate(transform.Rotation)
		worldVertex := transform.Position.Add(rotated)
		dot := float64(worldVertex.Dot(direction))
		if dot > maxDot {
			maxDot = dot
			maxVertex = worldVertex
		}
	}
	return maxVertex
}

type PolygonColliderComponentManager = ecs.ComponentManager[PolygonCollider]

func NewPolygonColliderComponentManager() PolygonColliderComponentManager {
	return ecs.NewComponentManager[PolygonCollider](PolygonColliderComponentId)
}

// ===========================
// Generic Collider
// ===========================

type GenericCollider struct {
	Shape      ColliderShape
	Layer      CollisionLayer
	Mask       CollisionMask
	Offset     vectors.Vec2
	AllowSleep bool
}

type GenericColliderComponentManager = ecs.ComponentManager[GenericCollider]

func NewGenericColliderComponentManager() GenericColliderComponentManager {
	return ecs.NewComponentManager[GenericCollider](GenericColliderComponentId)
}

type ColliderSleepState struct{}

type ColliderSleepStateComponentManager = ecs.ComponentManager[ColliderSleepState]

func NewColliderSleepStateComponentManager() ColliderSleepStateComponentManager {
	return ecs.NewComponentManager[ColliderSleepState](ColliderSleepStateComponentId)
}
