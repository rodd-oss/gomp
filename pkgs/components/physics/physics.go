package physics

import (
	"gomp_game/pkgs/engine"

	capnp "capnproto.org/go/capnp/v3"
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
)

type PhysicsController struct {
	body *cp.Body
}

var PhysicsComponent = engine.CreateNetworkComponent[PhysicsState](PhysicsState_TypeID, new(PhysicsController))

func (c *PhysicsController) Init() {
	c.body.SetPosition(cp.Vector{
		X: 50,
		Y: 75,
	})
}

func (c *PhysicsController) Update(dt float64) {
	pos := c.body.Position()

	pos.X++
	pos.Y++

	c.body.SetPosition(pos)
}

func (c *PhysicsController) CreateState(s *capnp.Segment) (PhysicsState, error) {
	return NewRootPhysicsState(s)
}

func (c *PhysicsController) OnStateRequest(state PhysicsState) PhysicsState {
	state.SetX(int32(c.body.Position().X))
	state.SetY(int32(c.body.Position().Y))

	return state
}

func (c *PhysicsController) OnStateUpdate(state PhysicsState) {
	pos := c.body.Position()

	pos.X = float64(state.X())
	pos.Y = float64(state.Y())

	c.body.SetPosition(pos)
}

func m() {
	w := ecs.NewWorld()
	PhysicsComponent.System.Each(w, func(e *ecs.Entry) {
		ctrl := PhysicsComponent.System.GetValue(e).Controller
		state := ctrl.OnStateRequest(PhysicsComponent.System.GetValue(e).State)
		state.Message().Marshal()
	})
}
