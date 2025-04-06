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

package main

import (
	"gomp"
	"gomp/examples/new-api/instances"
	"gomp/pkg/ecs"
	"gomp/stdsystems"
	"time"
)

func NewGame(initialScene gomp.AnyScene) Game {
	var game = Game{
		worlds:       make([]instances.World, 0),
		scenes:       make([]gomp.AnyScene, 0),
		lookup:       make(map[gomp.SceneId]int),
		renderSystem: stdsystems.NewRenderSystem(),
	}

	game.LoadScene(initialScene)
	game.SetActiveScene(initialScene.Id())

	return game
}

type Game struct {
	// Scene manager
	worlds []instances.World
	scenes []gomp.AnyScene
	lookup map[gomp.SceneId]int

	// Systems
	renderSystem stdsystems.RenderSystem

	// Utils
	shouldClose    bool
	currentSceneId gomp.SceneId
}

func (g *Game) Init() {
	g.renderSystem.Init()

	for i := range g.scenes {
		var scene = g.scenes[i]
		var world = &g.worlds[i]
		var systems = &world.Systems

		world.Init()

		systems.ColliderSystem.Init()
		systems.Velocity.Init()
		systems.CollisionDetectionBVH.Init()
		systems.CollisionResolution.Init()
		systems.AnimationSpriteMatrix.Init()
		systems.AnimationPlayer.Init()
		systems.SpriteMatrix.Init()
		systems.Sprite.Init()
		systems.YSort.Init()
		systems.Debug.Init()
		systems.AssetLib.Init()
		systems.Audio.Init()
		systems.SpatialAudio.Init()

		scene.Init(world)
	}
}

func (g *Game) Update(dt time.Duration) {
	for i := range g.scenes {
		g.scenes[i].Update(dt)
	}
}

func (g *Game) FixedUpdate(dt time.Duration) {
	for i := range g.worlds {
		var world = &g.worlds[i]
		var scene = g.scenes[i]
		var systems = &world.Systems

		systems.ColliderSystem.Run(dt)
		systems.Velocity.Run(dt)
		systems.CollisionDetectionBVH.Run(dt)
		systems.CollisionResolution.Run(dt)

		scene.FixedUpdate(dt)
	}
}

func (g *Game) Render(dt time.Duration) {
	var id = g.lookup[g.currentSceneId]
	var scene = g.scenes[id]
	var systems = &g.worlds[id].Systems

	systems.AnimationSpriteMatrix.Run()
	systems.AnimationPlayer.Run()
	systems.SpriteMatrix.Run()
	systems.Sprite.Run()
	systems.Debug.Run()
	systems.AssetLib.Run()
	systems.YSort.Run()
	systems.Audio.Run(dt)
	systems.SpatialAudio.Run(dt)

	g.shouldClose = g.renderSystem.Run(dt)

	scene.Render(dt)
}

func (g *Game) Destroy() {
	for i := range g.worlds {
		var scene = g.scenes[i]
		var world = &g.worlds[i]
		var systems = &world.Systems

		scene.Destroy()

		systems.Velocity.Destroy()
		systems.CollisionDetectionBVH.Destroy()
		systems.CollisionResolution.Destroy()
		systems.AnimationSpriteMatrix.Destroy()
		systems.AnimationPlayer.Destroy()
		systems.SpriteMatrix.Destroy()
		systems.Sprite.Destroy()
		systems.YSort.Destroy()
		systems.Debug.Destroy()
		systems.AssetLib.Destroy()
		systems.Audio.Destroy()
		systems.SpatialAudio.Destroy()

		world.Destroy()
	}

	g.renderSystem.Destroy()
}

func (g *Game) LoadScene(scene gomp.AnyScene) {
	g.scenes = append(g.scenes, scene)
	g.worlds = append(g.worlds, ecs.NewWorld(instances.NewComponentList(), instances.NewSystemList()))
	g.lookup[scene.Id()] = len(g.scenes) - 1
}

func (g *Game) SetActiveScene(id gomp.SceneId) {
	g.currentSceneId = id

	// Inject active scene world to render system
	var world = &g.worlds[g.lookup[id]]
	var components = &world.Components

	g.renderSystem.InjectWorld(
		&stdsystems.RenderInjector{
			EntityManager:                      &world.Entities,
			RlTexturePros:                      &components.RLTexturePro,
			Positions:                          &components.Position,
			Rotations:                          &components.Rotation,
			Scales:                             &components.Scale,
			AnimationPlayers:                   &components.AnimationPlayer,
			Tints:                              &components.Tint,
			Flips:                              &components.Flip,
			Renderables:                        &components.Renderable,
			AnimationStates:                    &components.AnimationState,
			Sprites:                            &components.Sprite,
			SpriteMatrixes:                     &components.SpriteMatrix,
			RenderOrders:                       &components.RenderOrder,
			BoxColliders:                       &components.ColliderBox,
			CircleColliders:                    &components.ColliderCircle,
			AABBs:                              &components.AABB,
			Collisions:                         &components.Collision,
			ColliderSleepStateComponentManager: &components.ColliderSleepState,
			BvhTrees:                           &components.BvhTree,
			Camera:                             &components.Cameras,
		})
}

func (g *Game) ShouldDestroy() bool {
	return g.shouldClose
}
