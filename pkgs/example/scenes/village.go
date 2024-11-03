package scenes

import (
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/entities"
	"gomp_game/pkgs/example/systems"
)

const VillageSceneName = "Village"

var VillageScene = engine.CreateScene(new(VillageSceneController))

type VillageSceneController struct {
	scene *engine.Scene
}

func (c *VillageSceneController) Load(scene *engine.Scene) (s []*engine.System, e []engine.Entity) {
	c.scene = scene

	s = append(s, &systems.PhysicsSystem)
	e = append(e, entities.Player)

	return s, e
}

func (c *VillageSceneController) Unload(scene *engine.Scene) {}
func (c *VillageSceneController) Update(dt float64)          {}
