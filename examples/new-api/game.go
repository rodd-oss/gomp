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
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/stdsystems"
	"time"
)

func NewGame(initialScene gomp.AnyScene) Game {
	var game = Game{
		worlds:          make([]instances.World, 0),
		scenes:          make([]gomp.AnyScene, 0),
		lookup:          make(map[gomp.SceneId]int),
		renderSystem:    stdsystems.NewRenderSystem(),
		osHandlerSystem: stdsystems.NewOSHandlerSystem(),
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
	renderSystem    stdsystems.RenderSystem
	osHandlerSystem stdsystems.OSHandlerSystem

	// Utils
	shouldClose    bool
	currentSceneId gomp.SceneId
}

func (g *Game) Init(engine *core.Engine) {
	g.osHandlerSystem.Init()
	g.renderSystem.Init()

	for i := range g.scenes {
		var scene = g.scenes[i]
		var world = &g.worlds[i]
		var systems = &world.Systems

		world.Init(engine)

		systems.ColliderSystem.Init()
		systems.Velocity.Init()
		systems.CollisionSetup.Init()
		systems.CollisionDetection.Init()
		//systems.CollisionDetectionBVH.Init()
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
		systems.AudioSettings.Init()
		systems.RenderCameras.Init()
		systems.Culling.Init()
		systems.TexturePositionSmooth.Init()

		scene.Init(world)
	}
}

func (g *Game) Update(dt time.Duration) {
	if g.osHandlerSystem.Run(dt) {
		g.shouldClose = true
		return
	}
	for i := range g.scenes {
		g.scenes[i].Update(dt)
	}
}

func (g *Game) FixedUpdate(dt time.Duration) {
	for i := range g.worlds {
		var world = &g.worlds[i]
		var scene = g.scenes[i]
		var systems = &world.Systems

		systems.Velocity.Run(dt)
		systems.ColliderSystem.Run(dt)
		systems.CollisionSetup.Run(dt)
		systems.CollisionDetection.Run(dt)
		//systems.CollisionDetectionBVH.Run(dt)
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
	systems.AudioSettings.Run(dt)

	scene.Render(dt)

	// Render all renderables with cameras
	systems.Culling.Run(dt)
	systems.RenderCameras.Run(dt)

	g.renderSystem.Run(dt)
}

func (g *Game) Destroy() {
	for i := range g.worlds {
		var scene = g.scenes[i]
		var world = &g.worlds[i]
		var systems = &world.Systems

		scene.Destroy()

		systems.Velocity.Destroy()
		systems.CollisionSetup.Destroy()
		systems.CollisionDetection.Destroy()
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
		systems.AudioSettings.Destroy()
		systems.RenderCameras.Destroy()
		systems.Culling.Destroy()
		systems.TexturePositionSmooth.Destroy()

		world.Destroy()
	}

	g.renderSystem.Destroy()
	g.osHandlerSystem.Destroy()
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
			EntityManager: &world.Entities,
			FrameBuffer2D: &components.FrameBuffer2D,
		})
}

func (g *Game) ShouldDestroy() bool {
	return g.shouldClose
}
