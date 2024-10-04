package components

import (
	"tomb_mates/internal/protos"

	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type PhysicsData struct {
	Body *cp.Body
}

func (p PhysicsData) Update(dt float64, e *ecs.Entry) error {
	if p.Body.IsSleeping() {
		return nil
	}

	pos := p.Body.Position()

	Transform.SetValue(e, TransformData{
		LocalPosition: math.NewVec2(pos.X, pos.Y),
	})

	unit := NetworkUnit.GetValue(e).Unit
	if unit != nil {
		unit.Position = &protos.Position{
			X: pos.X,
			Y: pos.Y,
		}
	}

	return nil
}

var Physics = ecs.NewComponentType[PhysicsData]()
