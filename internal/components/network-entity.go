package components

import (
	"tomb_mates/internal/protos"

	ecs "github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
)

type NetworkEntityData struct {
	Id        uint32
	Transform *protos.Transform
	Physics   *protos.Physics
	Skin      *protos.Skin
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
			ne.Transform.Position = iePosition
		}
		ieRotation := ieTransform.Rotation
		if ieRotation != nil {
			ne.Transform.Rotation = *ieRotation
		}
		ieScale := ieTransform.GetScale()
		if ieScale != nil {
			ne.Transform.Scale = ieScale
		}
	}

	iePhysics := patch.Physics
	if iePhysics != nil {
		ieVelocity := iePhysics.GetVelocity()
		if ieVelocity != nil {
			ne.Physics.Velocity = ieVelocity
		}
	}

	ieSkin := patch.Skin
	if ieSkin != nil {
		ne.Skin = ieSkin
	}
}

// Server-side
func (ne *NetworkEntityData) RequestPatch(entity *ecs.Entry) (patch *protos.PatchNetworkEntity) {
	transform := Transform.Get(entity)
	physics := Physics.Get(entity)

	patchTransform := ne.requestTransformPatch(transform)
	patchPhysics := ne.requestPhysicsPatch(physics)

	if patchTransform == nil && patchPhysics == nil {
		return nil
	}

	patch = &protos.PatchNetworkEntity{
		Transform: patchTransform,
		Physics:   patchPhysics,
	}

	return patch
}

func (ne *NetworkEntityData) requestTransformPatch(transform *transform.TransformData) (patchTransform *protos.PatchTransform) {
	if transform == nil {
		return nil
	}

	positionChanged := transform.LocalPosition.X != ne.Transform.Position.X || transform.LocalPosition.Y != ne.Transform.Position.Y
	rotationChanged := transform.LocalRotation != ne.Transform.Rotation
	scaleChanged := transform.LocalScale.X != ne.Transform.Scale.X || transform.LocalScale.Y != ne.Transform.Scale.Y

	if !(positionChanged || rotationChanged || scaleChanged) {
		return nil
	}

	patchTransform = &protos.PatchTransform{}

	if positionChanged {
		patchTransform.Position = &protos.Vector2{
			X: transform.LocalPosition.X,
			Y: transform.LocalPosition.Y,
		}
	}

	if rotationChanged {
		patchTransform.Rotation = &ne.Transform.Rotation
	}

	if scaleChanged {
		patchTransform.Scale = &protos.Vector2{
			X: transform.LocalScale.X,
			Y: transform.LocalScale.Y,
		}
	}

	return patchTransform
}

func (ne *NetworkEntityData) requestPhysicsPatch(physics *PhysicsData) (patchPhysics *protos.PatchPhysics) {
	if physics == nil {
		return nil
	}

	velocity := physics.Body.Velocity()

	velocityChanged := velocity.X != ne.Physics.Velocity.X || velocity.Y != ne.Physics.Velocity.Y

	if !(velocityChanged) {
		return
	}

	patchPhysics = &protos.PatchPhysics{}

	if velocityChanged {
		patchPhysics.Velocity = &protos.Vector2{
			X: velocity.X,
			Y: velocity.Y,
		}
	}

	return patchPhysics
}
