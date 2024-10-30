package scenes

import (
	"gomp_game/pkgs/engine"
	"log"
)

type VillageSceneController struct {
}

func (s *VillageSceneController) Update(dt float64) {
	log.Println("Village controller")
}

func (s *VillageSceneController) OnLoad(scene *engine.Scene) {
	log.Println("village scene load ")

	// scene.World.Create()
}

func (s *VillageSceneController) OnUnload(scene *engine.Scene) {
	// unload assets here
	log.Println("unload village scene")
}

var VillageScene = engine.NewScene(new(VillageSceneController))
