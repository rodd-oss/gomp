package components

import (
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type PhysicsData struct {
	Body *cp.Body
}

const (
	interpolationSpeed = 0.5
)

func (p PhysicsData) Update(dt float64, e *ecs.Entry, isClient bool) error {
	if p.Body.IsSleeping() {
		return nil
	}

	pos := p.Body.Position()
	posDelta := &math.Vec2{
		X: 0,
		Y: 0,
	}

	ne := NetworkEntity.GetValue(e)

	if ne.Transform != nil {
		networkPos := ne.Transform.Position
		if networkPos != nil {
			posDelta.X = pos.X - networkPos.X
			posDelta.Y = pos.Y - networkPos.Y
		}

		if isClient {
			p.Body.SetPosition(cp.Vector{
				X: pos.X - posDelta.X*interpolationSpeed,
				Y: pos.Y - posDelta.Y*interpolationSpeed,
			})
		}
	}

	Transform.SetValue(e, TransformData{
		LocalPosition: math.NewVec2(pos.X-posDelta.X/2, pos.Y-posDelta.Y/2),
	})

	// p.Body.SetVelocity(ne.Physics.Velocity.X, ne.Physics.Velocity.Y)

	return nil
}

var Physics = ecs.NewComponentType[PhysicsData]()
