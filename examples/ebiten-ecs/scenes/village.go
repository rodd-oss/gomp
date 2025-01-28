package scenes

import (
	"gomp/examples/ebiten-ecs/entities"
	"gomp/pkg/gomp"
)

var VillageScene = gomp.CreateScene("Village").
	AddEntities(
		entities.Hero(1),
	)
