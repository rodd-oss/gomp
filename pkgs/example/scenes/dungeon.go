package scenes

import (
	"gomp_game/pkgs/example/entities"
	"gomp_game/pkgs/gomp"
	"gomp_game/pkgs/gomp_utils/systems"
)

var DungeonScene = gomp.CreateScene("DungeonScene").
	AddEntities(entities.Player).
	AddSystems(systems.PhysicsSystem())
