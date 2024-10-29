package engine

type NetworkMode string

const (
	ServerMode NetworkMode = "server"
	ClientMode NetworkMode = "client"
)

type Network struct {
	Mode NetworkMode
}

func (n *Network) Update() {}

func (n *Network) GetAllPlayers() {}
