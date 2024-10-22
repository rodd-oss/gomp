package components

import (
	"math"

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
	pos.X = math.Round(pos.X)
	pos.Y = math.Round(pos.Y)
	// round body position to nearest integer
	p.Body.SetPosition(pos.Lerp(pos, interpolationSpeed))

	lastTransformPosition := Transform.GetValue(e).LocalPosition
	newTransformPosition := ecsmath.NewVec2(pos.X, pos.Y)

	posDelta := &ecsmath.Vec2{
		X: p.Body.Position().X - lastTransformPosition.X,
		Y: p.Body.Position().Y - lastTransformPosition.Y,
	}

	if isClient {
		Transform.SetValue(e, TransformData{
			LocalPosition: newTransformPosition.Add(posDelta.MulScalar(-0.66)),
		})
	} else {
		Transform.SetValue(e, TransformData{
			LocalPosition: newTransformPosition,
		})
	}

	return nil
}

var Physics = ecs.NewComponentType[PhysicsData]()
