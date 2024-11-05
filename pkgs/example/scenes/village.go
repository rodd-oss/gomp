package scenes

import (
	"gomp_game/pkgs/example/entities"
	"gomp_game/pkgs/gomp"
	"gomp_game/pkgs/gomp_utils/systems"
)

var VillageScene = gomp.CreateScene("Village").
	AddEntities(
		entities.Player,
		entities.Player,
	).
	AddSystems(
		systems.PhysicsSystem(),
	)
