package components

import (
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	ecsmath "github.com/yohamta/donburi/features/math"
)

type PhysicsData struct {
	Body *cp.Body
}

const (
	interpolationSpeed = 0.25
)

func (p PhysicsData) Update(dt float64, e *ecs.Entry, isClient bool) error {
	if p.Body.IsSleeping() {
		return nil
	}

	pos := p.Body.Position()
	posDelta := &ecsmath.Vec2{
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

		// Force sync client body pos to server body pos
		if isClient {
			if p.Body.Velocity().X == 0 && p.Body.Velocity().Y == 0 {
				p.Body.SetPosition(cp.Vector{
					X: pos.X - posDelta.X*interpolationSpeed,
					Y: pos.Y - posDelta.Y*interpolationSpeed,
				})
			}
		}
	}

	if isClient {
		Transform.SetValue(e, TransformData{
			LocalPosition: ecsmath.NewVec2(
				p.Body.Position().X,
				p.Body.Position().Y),
		})

		p.Body.SetVelocity(ne.Physics.Velocity.X, ne.Physics.Velocity.Y)
	} else {
		Transform.SetValue(e, TransformData{
			LocalPosition: ecsmath.NewVec2(p.Body.Position().X, p.Body.Position().Y),
		})
	}

	return nil
}

var Physics = ecs.NewComponentType[PhysicsData]()
