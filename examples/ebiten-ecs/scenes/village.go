package scenes

import (
	"gomp/examples/ebiten-ecs/entities"
	"gomp/pkgs/gomp"
)

var VillageScene = gomp.CreateScene("Village").
	AddEntities(
		entities.Hero(1),
	)
