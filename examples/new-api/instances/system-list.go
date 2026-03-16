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
	"gomp"
	"gomp/examples/new-api/assets"
	"gomp/examples/new-api/systems"
	"gomp/stdsystems"
)

func NewSystemList() SystemList {
	newSystemList := SystemList{
		Debug:                    stdsystems.NewDebugSystem(),
		Velocity:                 stdsystems.NewVelocitySystem(),
		Network:                  stdsystems.NewNetworkSystem(),
		NetworkReceive:           stdsystems.NewNetworkReceiveSystem(),
		NetworkSend:              stdsystems.NewNetworkSendSystem(),
		AnimationSpriteMatrix:    stdsystems.NewAnimationSpriteMatrixSystem(),
		AnimationPlayer:          stdsystems.NewAnimationPlayerSystem(),
		TextureRenderSpriteSheet: stdsystems.NewTextureRenderSpriteSheetSystem(),
		Sprite:                   stdsystems.NewSpriteSystem(),
		SpriteMatrix:             stdsystems.NewSpriteMatrixSystem(),
		AssetLib:                 stdsystems.NewAssetLibSystem([]gomp.AnyAssetLibrary{&assets.Textures, &assets.Audio}),
		YSort:                    stdsystems.NewYSortSystem(),
		CollisionDetectionGrid:   stdsystems.NewCollisionDetectionGridSystem(),
		CollisionDetectionBVH:    stdsystems.NewCollisionDetectionBVHSystem(),
		ColliderSystem:           stdsystems.NewColliderSystem(),
		CollisionSetup:           stdsystems.NewCollisionSetupSystem(),
		CollisionDetection:       stdsystems.NewCollisionDetectionSystem(),
		CollisionResolution:      stdsystems.NewCollisionResolutionSystem(),
		TexturePositionSmooth:    stdsystems.NewTexturePositionSmoothSystem(),
		RenderCameras:            stdsystems.NewRender2DCamerasSystem(),
		Culling:                  stdsystems.NewCullingSystem(),
		Player:                   systems.NewPlayerSystem(),
		RenderBogdan:             systems.NewRenderBogdanSystem(),
		Audio:                    systems.NewAudioSystem(),
		AudioSettings:            systems.NewAudioSettingsSystem(),
		SpatialAudio:             systems.NewSpatialAudioSystem(),
		DampingSystem:            systems.NewDampingSystem(),
		AssteroddSystem:          systems.NewAssteroddSystem(),
		CollisionHandler:         systems.NewCollisionHandlerSystem(),
		SpaceshipIntents:         systems.NewSpaceshipIntentsSystem(),
		SpaceSpawner:             systems.NewSpaceSpawnerSystem(),
		Hp:                       systems.NewHpSystem(),
		MainCamera:               systems.NewMainCameraSystem(),
		TextureRect:              systems.NewTextureRectSystem(),
		TextureCircle:            systems.NewTextureCircleSystem(),
		Minimap:                  systems.NewMinimapSystem(),
		DebugInfo:                systems.NewDebugInfoSystem(),
		RenderOverlay:            systems.NewRenderOverlaySystem(),
	}

	return newSystemList
}

type SystemList struct {
	Debug                    stdsystems.DebugSystem
	Velocity                 stdsystems.VelocitySystem
	Network                  stdsystems.NetworkSystem
	NetworkReceive           stdsystems.NetworkReceiveSystem
	NetworkSend              stdsystems.NetworkSendSystem
	AnimationSpriteMatrix    stdsystems.AnimationSpriteMatrixSystem
	AnimationPlayer          stdsystems.AnimationPlayerSystem
	TextureRenderSpriteSheet stdsystems.TextureRenderSpriteSheetSystem
	Sprite                   stdsystems.SpriteSystem
	SpriteMatrix             stdsystems.SpriteMatrixSystem
	AssetLib                 stdsystems.AssetLibSystem
	YSort                    stdsystems.YSortSystem
	CollisionDetectionGrid   stdsystems.CollisionDetectionGridSystem
	CollisionDetectionBVH    stdsystems.CollisionDetectionBVHSystem
	ColliderSystem           stdsystems.ColliderSystem
	CollisionSetup           stdsystems.CollisionSetupSystem
	CollisionDetection       stdsystems.CollisionDetectionSystem
	CollisionResolution      stdsystems.CollisionResolutionSystem
	TexturePositionSmooth    stdsystems.TexturePositionSmoothSystem
	RenderCameras            stdsystems.Render2DCamerasSystem
	Culling                  stdsystems.CullingSystem
	RenderBogdan             systems.RenderBogdanSystem
	Player                   systems.PlayerSystem
	Audio                    systems.AudioSystem
	AudioSettings            systems.AudioSettingsSystem
	SpatialAudio             systems.SpatialAudioSystem
	DampingSystem            systems.DampingSystem
	AssteroddSystem          systems.AssteroddSystem
	CollisionHandler         systems.CollisionHandlerSystem
	SpaceshipIntents         systems.SpaceshipIntentsSystem
	SpaceSpawner             systems.SpaceSpawnerSystem
	Hp                       systems.HpSystem
	MainCamera               systems.MainCameraSystem
	Minimap                  systems.MinimapSystem
	DebugInfo                systems.DebugInfoSystem
	TextureRect              systems.TextureRectSystem
	TextureCircle            systems.TextureCircleSystem
	RenderOverlay            systems.RenderOverlaySystem
}
