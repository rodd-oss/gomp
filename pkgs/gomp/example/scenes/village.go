package scenes

import (
	"gomp_game/pkgs/gomp"
	"gomp_game/pkgs/gomp/example/entities"
	"gomp_game/pkgs/gomp/example/systems"
)

var VillageScene = gomp.CreateScene("Village").
	AddEntities(
		entities.Player,
		entities.Player,
	).
	AddSystems(
		systems.PhysicsSystem(),
	)
