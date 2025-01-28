package scenes

import (
	"gomp/examples/ebiten-ecs/entities"
	"gomp/pkg/gomp"
)

var DungeonScene = gomp.CreateScene("DungeonScene").
	AddEntities(entities.Hero(1))
