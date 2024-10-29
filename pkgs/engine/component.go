package engine

type State[T any, S any] interface {
	Get() *T
	GetSerialized() *S

	Set(t *T) error

	Serialize()
}

type NetworkSyncComponent[T any, S any] interface {
	Update(dt float64)

	State() State[T, S]
}

func CreateNetworkComponent[T, S any]() NetworkSyncComponent[T, S] {}

type Component[T any, S any] interface {
	Update(dt float64)
}

// used by devs
type pcs struct {
	X, Y float64
}

var physicsComponent = CreateNetworkComponent[pcs, []byte]()

func m() {
	physicsComponent.State().GetSerialized()
}
