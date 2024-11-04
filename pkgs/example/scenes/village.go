package scenes

import (
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/entities"
	"gomp_game/pkgs/example/systems"
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
