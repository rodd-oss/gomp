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

package instances

import (
	"gomp/examples/new-api/components"
	"gomp/stdcomponents"
)

type ComponentList struct {
	Position           stdcomponents.PositionComponentManager
	Rotation           stdcomponents.RotationComponentManager
	Scale              stdcomponents.ScaleComponentManager
	Velocity           stdcomponents.VelocityComponentManager
	Flip               stdcomponents.FlipComponentManager
	Sprite             stdcomponents.SpriteComponentManager
	SpriteMatrix       stdcomponents.SpriteMatrixComponentManager
	Tint               stdcomponents.TintComponentManager
	AnimationPlayer    stdcomponents.AnimationPlayerComponentManager
	AnimationState     stdcomponents.AnimationStateComponentManager
	RLTexturePro       stdcomponents.RLTextureProComponentManager
	Network            stdcomponents.NetworkComponentManager
	Renderable         stdcomponents.RenderableComponentManager
	YSort              stdcomponents.YSortComponentManager
	RenderOrder        stdcomponents.RenderOrderComponentManager
	GenericCollider    stdcomponents.GenericColliderComponentManager
	ColliderBox        stdcomponents.BoxColliderComponentManager
	ColliderCircle     stdcomponents.CircleColliderComponentManager
	ColliderSleepState stdcomponents.ColliderSleepStateComponentManager
	Collision          stdcomponents.CollisionComponentManager
	AABB               stdcomponents.AABBComponentManager
	SpatialIndex       stdcomponents.SpatialIndexComponentManager
	RigidBody          stdcomponents.RigidBodyComponentManager
	BvhTree            stdcomponents.BvhTreeComponentManager

	Health               components.HpComponentManager
	Controller           components.ControllerComponentManager
	PlayerTag            components.PlayerTagComponentManager
	BulletTag            components.BulletTagComponentManager
	AsteroidTag          components.AsteroidComponentManager
	SpaceSpawnerTag      components.SpaceSpawnerComponentManager
	Wall                 components.WallTagComponentManager
	Weapon               components.WeaponComponentManager
	SpaceshipIntent      components.SpaceshipIntentComponentManager
	AsteroidSceneManager components.AsteroidSceneManagerComponentManager
	SoundEffects         components.SoundEffectsComponentManager
	MainCamera           stdcomponents.MainCameraComponentManager
	PipCamera            stdcomponents.PipCameraComponentManager
	MinimapCamera        stdcomponents.MinimapCameraComponentManager
	RenderTexture2D      stdcomponents.RenderTexture2DComponentManager
}

func NewComponentList() ComponentList {
	return ComponentList{
		Position:           stdcomponents.NewPositionComponentManager(),
		Rotation:           stdcomponents.NewRotationComponentManager(),
		Scale:              stdcomponents.NewScaleComponentManager(),
		Velocity:           stdcomponents.NewVelocityComponentManager(),
		Flip:               stdcomponents.NewFlipComponentManager(),
		Sprite:             stdcomponents.NewSpriteComponentManager(),
		SpriteMatrix:       stdcomponents.NewSpriteMatrixComponentManager(),
		Tint:               stdcomponents.NewTintComponentManager(),
		AnimationPlayer:    stdcomponents.NewAnimationPlayerComponentManager(),
		AnimationState:     stdcomponents.NewAnimationStateComponentManager(),
		RLTexturePro:       stdcomponents.NewRlTextureProComponentManager(),
		Network:            stdcomponents.NewNetworkComponentManager(),
		Renderable:         stdcomponents.NewRenderableComponentManager(),
		YSort:              stdcomponents.NewYSortComponentManager(),
		RenderOrder:        stdcomponents.NewRenderOrderComponentManager(),
		GenericCollider:    stdcomponents.NewGenericColliderComponentManager(),
		ColliderBox:        stdcomponents.NewBoxColliderComponentManager(),
		ColliderCircle:     stdcomponents.NewCircleColliderComponentManager(),
		ColliderSleepState: stdcomponents.NewColliderSleepStateComponentManager(),
		Collision:          stdcomponents.NewCollisionComponentManager(),
		AABB:               stdcomponents.NewAABBComponentManager(),
		SpatialIndex:       stdcomponents.NewSpatialIndexComponentManager(),
		RigidBody:          stdcomponents.NewRigidBodyComponentManager(),
		BvhTree:            stdcomponents.NewBvhTreeComponentManager(),

		Health:               components.NewHealthComponentManager(),
		Controller:           components.NewControllerComponentManager(),
		PlayerTag:            components.NewPlayerTagComponentManager(),
		BulletTag:            components.NewBulletTagComponentManager(),
		Wall:                 components.NewWallComponentManager(),
		AsteroidTag:          components.NewAsteroidTagComponentManager(),
		SpaceSpawnerTag:      components.NewSpaceSpawnerTagComponentManager(),
		Weapon:               components.NewWeaponComponentManager(),
		SpaceshipIntent:      components.NewSpaceshipIntentComponentManager(),
		AsteroidSceneManager: components.NewAsteroidSceneManagerComponentManager(),
		SoundEffects:         components.NewSoundEffectsComponentManager(),
		MainCamera:           stdcomponents.NewMainCameraComponentManager(),
		PipCamera:            stdcomponents.NewPipCameraComponentManager(),
		MinimapCamera:        stdcomponents.NewMinimapCameraComponentManager(),
		RenderTexture2D:      stdcomponents.NewRenderTexture2DComponentManager(),
	}
}
