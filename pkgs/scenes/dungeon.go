package scenes

import (
	"gomp_game/pkgs/engine"
	"log"
)

const DungeonSceneName = "Dungeon"

type DungeonSceneController struct {
}

func (s *DungeonSceneController) Update(dt float64) {
	log.Println("Dungeon contoller")
}

func (s *DungeonSceneController) OnLoad(scene *engine.Scene) {
	log.Println("Dungeon scene load")

	// scene.World.Create()
}

func (s *DungeonSceneController) OnUnload(scene *engine.Scene) {
	// unload assets here
	log.Println("unload Dungeon scene")
}

var DungeonScene = engine.NewScene(new(DungeonSceneController))
