/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- HromRu Donated 2 500 RUB
<- kotyamatroskin Donated 400 RUB
<- r_uslashk_a Donated 100 RUB
<- mitwelve Donated 100 RUB

Thank you for your support!
*/

package scenes

import (
	"gomp"
	"gomp/examples/new-api/instances"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"time"
)

/*
TITLE:
- Roddsteroids
- Asteroidds
- Asteroids
+ Assterodds

USER FLOW:
- press space to start
+ spaceship AD - turn ; WS - gas
+ spaceship HP = 3
+ asteroids coming from the top
+ asteroids HP=3-9
- random bonus from asteroid (hp, weapon, move speed)
- asteroid spawn rate increases over the time
- player scores (optional save scores over the session)

ENTITIES and COMPONENTS:
- scene-manager: sceneState, playerScore, playerHighScore
+ spaceship: position, rotation, scale, hp, velocity, sprite, collider, playerTag, weapon
+ asteroid: position, rotation, scale, hp, velocity, sprite, collider, asteroidTag, scoreReward
+ asteroid-spawner: position, velocity, asteroidSpawnerTag
+ bullet: position, rotation, scale, velocity, sprite, collider, damage, bulletTag
- bonus: position, rotation, scale, velocity, sprite, collider, bonusTag

COMPONENTS:
- position, rotation, scale, velocity
- hp, damage, scoreReward, weapon
- asteroidTag, bulletTag, bonusTag, playerTag, asteroidSpawnerTag
- sprite, collider
*/

func NewAssteroddScene() AssteroddScene {
	return AssteroddScene{}
}

type AssteroddScene struct {
	Game  core.AnyGame
	World *instances.World
}

func (s *AssteroddScene) Init(world ecs.AnyWorld) {
	s.World = world.(*instances.World)

	s.World.Systems.DampingSystem.Init()
	s.World.Systems.SpaceSpawner.Init()
	s.World.Systems.AssteroddSystem.Init()
	s.World.Systems.CollisionHandler.Init()
	s.World.Systems.SpaceshipIntents.Init()
	s.World.Systems.MainCamera.Init()
	s.World.Systems.Minimap.Init()
	s.World.Systems.TexturePositionSmooth.Init()
	s.World.Systems.TextureRect.Init()
	s.World.Systems.TextureCircle.Init()
	s.World.Systems.DebugInfo.Init()
	s.World.Systems.RenderOverlay.Init()

}

func (s *AssteroddScene) Update(dt time.Duration) gomp.SceneId {
	s.World.Systems.AssteroddSystem.Run(dt)

	return AssteroddSceneId
}

func (s *AssteroddScene) FixedUpdate(dt time.Duration) {
	s.World.Systems.SpaceshipIntents.Run(dt)
	s.World.Systems.DampingSystem.Run(dt)
	//s.World.Systems.SpaceSpawner.Run(dt)
	s.World.Systems.CollisionHandler.Run(dt)
	s.World.Systems.Hp.Run(dt)
}

func (s *AssteroddScene) Render(dt time.Duration) {
	// Camera game logic flow
	s.World.Systems.MainCamera.Run(dt)
	s.World.Systems.Minimap.Run(dt)

	s.World.Systems.DebugInfo.Run(dt)

	// Optimized primitives
	s.World.Systems.TextureRect.Run(dt)
	s.World.Systems.TextureCircle.Run(dt)

	// Over cameras render example
	s.World.Systems.RenderOverlay.Run(dt)
}

func (s *AssteroddScene) Destroy() {
	s.World.Systems.DampingSystem.Destroy()
	s.World.Systems.SpaceSpawner.Destroy()
	s.World.Systems.AssteroddSystem.Destroy()
	s.World.Systems.CollisionHandler.Destroy()
	s.World.Systems.SpaceshipIntents.Destroy()
	s.World.Systems.MainCamera.Destroy()
	s.World.Systems.TexturePositionSmooth.Destroy()
	s.World.Systems.TextureRect.Destroy()
	s.World.Systems.TextureCircle.Destroy()
	s.World.Systems.Minimap.Destroy()
	s.World.Systems.DebugInfo.Destroy()
}

func (s *AssteroddScene) OnEnter() {

}

func (s *AssteroddScene) OnExit() {

}

func (s *AssteroddScene) Id() gomp.SceneId {
	return AssteroddSceneId
}

var _ gomp.AnyScene = (*AssteroddScene)(nil)
