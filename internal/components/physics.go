package components

import (
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
)

type PhysicsData struct {
	Body *cp.Body
}

var Physics = ecs.NewComponentType[PhysicsData]()
