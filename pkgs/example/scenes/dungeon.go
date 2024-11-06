package scenes

import (
	"gomp_game/pkgs/example/entities"
	"gomp_game/pkgs/gomp"
)

var DungeonScene = gomp.CreateScene("DungeonScene").
	AddEntities(entities.Player())
