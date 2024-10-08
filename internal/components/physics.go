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
	posDelta := &math.Vec2{
		X: 0,
		Y: 0,
	}

	networkEntityTransform := NetworkEntity.GetValue(e).Transform
	if networkEntityTransform != nil {
		networkPos := NetworkEntity.GetValue(e).Transform.Position
		if networkPos != nil {
			posDelta.X = pos.X - networkPos.X
			posDelta.Y = pos.Y - networkPos.Y
		}
	}

	Transform.SetValue(e, TransformData{
		LocalPosition: math.NewVec2(pos.X-posDelta.X/2, pos.Y-posDelta.Y/2),
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
