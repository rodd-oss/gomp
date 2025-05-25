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
	Position              stdcomponents.PositionComponentManager
	Rotation              stdcomponents.RotationComponentManager
	Scale                 stdcomponents.ScaleComponentManager
	Velocity              stdcomponents.VelocityComponentManager
	Flip                  stdcomponents.FlipComponentManager
	Sprite                stdcomponents.SpriteComponentManager
	SpriteMatrix          stdcomponents.SpriteMatrixComponentManager
	Tint                  stdcomponents.TintComponentManager
	AnimationPlayer       stdcomponents.AnimationPlayerComponentManager
	AnimationState        stdcomponents.AnimationStateComponentManager
	RLTexturePro          stdcomponents.RLTextureProComponentManager
	Network               stdcomponents.NetworkComponentManager
	Renderable            stdcomponents.RenderableComponentManager
	YSort                 stdcomponents.YSortComponentManager
	RenderOrder           stdcomponents.RenderOrderComponentManager
	GenericCollider       stdcomponents.GenericColliderComponentManager
	ColliderBox           stdcomponents.BoxColliderComponentManager
	ColliderCircle        stdcomponents.CircleColliderComponentManager
	ColliderSleepState    stdcomponents.ColliderSleepStateComponentManager
	Collision             stdcomponents.CollisionComponentManager
	AABB                  stdcomponents.AABBComponentManager
	SpatialIndex          stdcomponents.SpatialHashComponentManager
	RigidBody             stdcomponents.RigidBodyComponentManager
	BvhTree               stdcomponents.BvhTreeComponentManager
	Cameras               stdcomponents.CameraComponentManager
	TexturePositionSmooth stdcomponents.TexturePositionSmoothComponentManager
	FrameBuffer2D         stdcomponents.FrameBuffer2DComponentManager
	CollisionGrid         stdcomponents.CollisionGridComponentManager
	CollisionChunk        stdcomponents.CollisionChunkComponentManager
	CollisionCell         stdcomponents.CollisionCellComponentManager
	CollisionGridMember   stdcomponents.CollisionGridMemberComponentManager

	Health               components.HpComponentManager
	Controller           components.ControllerComponentManager
	PlayerTag            components.PlayerTagComponentManager
	BulletTag            components.BulletTagComponentManager
	AsteroidTag          components.AsteroidTagComponentManager
	SpaceSpawnerTag      components.SpaceSpawnerTagComponentManager
	Wall                 components.WallComponentManager
	Weapon               components.WeaponComponentManager
	SpaceshipIntent      components.SpaceshipIntentComponentManager
	AsteroidSceneManager components.AsteroidSceneManagerComponentManager
	SoundEffects         components.SoundEffectComponentManager
	SpatialAudio         components.SpatialAudioComponentManager
	TextureRect          components.TextureRectComponentManager
	PrimitiveCircle      components.TextureCircleComponentManager
	RenderVisible        stdcomponents.RenderVisibleComponentManager
}

func NewComponentList() ComponentList {
	return ComponentList{
		Position:              stdcomponents.NewPositionComponentManager(),
		Rotation:              stdcomponents.NewRotationComponentManager(),
		Scale:                 stdcomponents.NewScaleComponentManager(),
		Velocity:              stdcomponents.NewVelocityComponentManager(),
		Flip:                  stdcomponents.NewFlipComponentManager(),
		Sprite:                stdcomponents.NewSpriteComponentManager(),
		SpriteMatrix:          stdcomponents.NewSpriteMatrixComponentManager(),
		Tint:                  stdcomponents.NewTintComponentManager(),
		AnimationPlayer:       stdcomponents.NewAnimationPlayerComponentManager(),
		AnimationState:        stdcomponents.NewAnimationStateComponentManager(),
		RLTexturePro:          stdcomponents.NewRLTextureProComponentManager(),
		Network:               stdcomponents.NewNetworkComponentManager(),
		Renderable:            stdcomponents.NewRenderableComponentManager(),
		RenderVisible:         stdcomponents.NewRenderVisibleComponentManager(),
		YSort:                 stdcomponents.NewYSortComponentManager(),
		RenderOrder:           stdcomponents.NewRenderOrderComponentManager(),
		GenericCollider:       stdcomponents.NewGenericColliderComponentManager(),
		ColliderBox:           stdcomponents.NewBoxColliderComponentManager(),
		ColliderCircle:        stdcomponents.NewCircleColliderComponentManager(),
		ColliderSleepState:    stdcomponents.NewColliderSleepStateComponentManager(),
		Collision:             stdcomponents.NewCollisionComponentManager(),
		AABB:                  stdcomponents.NewAABBComponentManager(),
		SpatialIndex:          stdcomponents.NewSpatialHashComponentManager(),
		RigidBody:             stdcomponents.NewRigidBodyComponentManager(),
		BvhTree:               stdcomponents.NewBvhTreeComponentManager(),
		Cameras:               stdcomponents.NewCameraComponentManager(),
		FrameBuffer2D:         stdcomponents.NewFrameBuffer2DComponentManager(),
		TexturePositionSmooth: stdcomponents.NewTexturePositionSmoothComponentManager(),
		CollisionGrid:         stdcomponents.NewCollisionGridComponentManager(),
		CollisionChunk:        stdcomponents.NewCollisionChunkComponentManager(),
		CollisionCell:         stdcomponents.NewCollisionCellComponentManager(),
		CollisionGridMember:   stdcomponents.NewCollisionGridMemberComponentManager(),

		Health:               components.NewHpComponentManager(),
		Controller:           components.NewControllerComponentManager(),
		PlayerTag:            components.NewPlayerTagComponentManager(),
		BulletTag:            components.NewBulletTagComponentManager(),
		Wall:                 components.NewWallComponentManager(),
		AsteroidTag:          components.NewAsteroidTagComponentManager(),
		SpaceSpawnerTag:      components.NewSpaceSpawnerTagComponentManager(),
		Weapon:               components.NewWeaponComponentManager(),
		SpaceshipIntent:      components.NewSpaceshipIntentComponentManager(),
		AsteroidSceneManager: components.NewAsteroidSceneManagerComponentManager(),
		SoundEffects:         components.NewSoundEffectComponentManager(),
		SpatialAudio:         components.NewSpatialAudioComponentManager(),
		TextureRect:          components.NewTextureRectComponentManager(),
		PrimitiveCircle:      components.NewTextureCircleComponentManager(),
	}
}
