package scenes

import (
	"gomp_game/pkgs/gomp-v1/engine"
	"gomp_game/pkgs/gomp-v1/example/entities"
	"gomp_game/pkgs/gomp-v1/example/systems"
)

var VillageScene = engine.CreateScene("Village").
	AddEntities(
		entities.Player,
		entities.Player,
	).
	AddSystems(
		systems.PhysicsSystem(),
	).
	AddRenderSystems(
		systems.RenderSystem(),
	)
