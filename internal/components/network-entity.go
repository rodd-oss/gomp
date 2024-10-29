/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package components

import (
	"gomp_game/internal/protos"

	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type NetworkEntityData struct {
	Id uint32

	LastPatch *protos.PatchNetworkEntity

	Body      *cp.Body
	Transform *TransformData
}

var NetworkEntity = ecs.NewComponentType[*NetworkEntityData]()

// Client-side
func (ne *NetworkEntityData) ApplyPatch(patch *protos.PatchNetworkEntity) {
	if patch == nil {
		return
	}

	ieTransform := patch.Transform
	if ieTransform != nil {
		iePosition := ieTransform.GetPosition()
		if iePosition != nil {
			ne.Body.SetPosition(cp.Vector{X: iePosition.X, Y: iePosition.Y})
		}
		ieRotation := ieTransform.Rotation
		if ieRotation != nil {
			ne.Transform.LocalRotation = *ieRotation
		}
		ieScale := ieTransform.GetScale()
		if ieScale != nil {
			ne.Transform.LocalScale = math.Vec2{X: ieScale.X, Y: ieScale.Y}
		}
	}

	iePhysics := patch.Physics
	if iePhysics != nil {
		ieVelocity := iePhysics.GetVelocity()
		if ieVelocity != nil {
			ne.Body.SetVelocity(ieVelocity.X, ieVelocity.Y)
		}
		iePosition := iePhysics.GetPosition()
		if iePosition != nil {
			ne.Body.SetPosition(cp.Vector{X: iePosition.X, Y: iePosition.Y})
		}
	}
}

// Server-side
func (ne *NetworkEntityData) RequestPatch(entity *ecs.Entry) (patch *protos.PatchNetworkEntity) {

	patchTransform := ne.requestTransformPatch()
	patchPhysics := ne.requestPhysicsPatch()

	if patchPhysics == nil {
		return nil
	}

	patch = &protos.PatchNetworkEntity{
		Transform: patchTransform,
		Physics:   patchPhysics,
	}

	ne.LastPatch = patch

	return patch
}

func (ne *NetworkEntityData) requestTransformPatch() (patchTransform *protos.PatchTransform) {
	transform := ne.Transform

	if transform == nil {
		return nil
	}

	var positionChanged, rotationChanged, scaleChanged bool = false, false, false

	if ne.LastPatch != nil {
		if ne.LastPatch.Transform != nil {
			if ne.LastPatch.Transform.Position != nil {
				positionChanged = transform.LocalPosition.X != ne.LastPatch.Transform.Position.X || transform.LocalPosition.Y != ne.LastPatch.Transform.Position.Y
			}

			if ne.LastPatch.Transform.Rotation != nil {
				rotationChanged = transform.LocalRotation != *ne.LastPatch.Transform.Rotation
			}

			if ne.LastPatch.Transform.Scale != nil {
				scaleChanged = transform.LocalScale.X != ne.LastPatch.Transform.Scale.X || transform.LocalScale.Y != ne.LastPatch.Transform.Scale.Y
			}
		}

		if !(positionChanged || rotationChanged || scaleChanged) {
			return nil
		}
	}

	patchTransform = &protos.PatchTransform{}

	if positionChanged {
		patchTransform.Position = &protos.Vector2{
			X: transform.LocalPosition.X,
			Y: transform.LocalPosition.Y,
		}
	}

	if rotationChanged {
		patchTransform.Rotation = ne.LastPatch.Transform.Rotation
	}

	if scaleChanged {
		patchTransform.Scale = &protos.Vector2{
			X: transform.LocalScale.X,
			Y: transform.LocalScale.Y,
		}
	}

	return patchTransform
}

func (ne *NetworkEntityData) requestPhysicsPatch() (patchPhysics *protos.PatchPhysics) {
	if ne.Body == nil {
		return nil
	}

	velocity := ne.Body.Velocity()
	position := ne.Body.Position()

	if ne.LastPatch != nil {

		velocityChanged := velocity.X != ne.LastPatch.Physics.Velocity.X || velocity.Y != ne.LastPatch.Physics.Velocity.Y

		if !(velocityChanged) {
			return nil
		}
	}

	patchPhysics = &protos.PatchPhysics{}

	patchPhysics.Velocity = &protos.Vector2{
		X: velocity.X,
		Y: velocity.Y,
	}

	patchPhysics.Position = &protos.Vector2{
		X: position.X,
		Y: position.Y,
	}

	return patchPhysics
}
