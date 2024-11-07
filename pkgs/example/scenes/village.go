package scenes

import (
	"gomp_game/pkgs/example/entities"
	"gomp_game/pkgs/gomp"
)

var VillageScene = gomp.CreateScene("Village").
	AddEntities(
		entities.Enemy,
		entities.Enemy2,
	)
