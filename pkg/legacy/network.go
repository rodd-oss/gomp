/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package legacy

type NetworkMode int

const (
	NetworkMode_None NetworkMode = iota
	NetworkMode_Host
	NetworkMode_Client
)

type NetworkPlayerId int
type NetworkPlayer struct{}

type HostClient struct {
	id        NetworkPlayerId
	transport HostClientTransport
}

type HostClientTransport interface {
	Send(msg []byte)
	Close() error
}

type Network struct {
	Mode NetworkMode
	// Players map[NetworkPlayerId]NetworkPlayer

	Host   NetworkHost
	Client NetworkClient

	NewPlayers     chan NetworkPlayer
	RemovedPlayers chan int
	InEvents       chan []byte
	OutEvents      chan []byte
}

func (n *Network) Update(dt float64) {
	switch n.Mode {
	case NetworkMode_None:
		break
	case NetworkMode_Host:
		n.UpdateHost(dt)
	case NetworkMode_Client:
		n.UpdateClient(dt)
	default:
		panic("invalid network mode")
	}
}

func (n *Network) UpdateHost(dt float64) {
	newEventsLen := len(n.OutEvents)
	for i := 0; i < newEventsLen; i++ {
		event := <-n.OutEvents
		n.Host.Broadcast(event)
	}

	newClientsLen := len(n.Host.register)
	for i := 0; i < newClientsLen; i++ {
		client := <-n.Host.register
		n.Host.clients[client.id] = client
		// TODO: create player as a world entity
		// n.Players[ebiten-client.id] = NetworkPlayer{}

		if i >= 9 {
			break
		}
	}

	removedClientsLen := len(n.Host.unregister)
	for i := 0; i < removedClientsLen; i++ {
		id := <-n.Host.unregister
		delete(n.Host.clients, id)
		// delete(n.Players, id)

		if i >= 9 {
			break
		}
	}

	newMessagesLen := len(n.Host.receiver)
	for i := 0; i < newMessagesLen; i++ {
		msg := <-n.Host.receiver
		n.InEvents <- msg
	}
}

func (n *Network) UpdateClient(dt float64) {
	newEventsLen := len(n.OutEvents)
	for i := 0; i < newEventsLen; i++ {
		event := <-n.OutEvents
		n.Client.Send(event)
	}

	msgs := n.Client.transport.Receive()
	msgLen := len(msgs)
	for i := 0; i < msgLen; i++ {
		// check if msg is AddPlayer or RemovePlayer
		n.InEvents <- msgs[i]
	}
}

type ClientTransport interface {
	Connect() (NetworkPlayerId, error)
	Close()
	Send([]byte)
	Receive() [][]byte
}

type NetworkClient struct {
	id        NetworkPlayerId
	transport ClientTransport
}

func (c *NetworkClient) Connect(transport ClientTransport) (err error) {
	c.id, err = c.transport.Connect()
	if err != nil {
		return err
	}

	c.transport = transport

	return nil
}

func (c *NetworkClient) Close() {
	c.transport.Close()
}

func (c *NetworkClient) Send(msg []byte) {
	c.transport.Send(msg)
}

func (c *NetworkClient) Receive() {

}

type NetworkHost struct {
	receiver   chan []byte
	MaxClients int
	clients    map[NetworkPlayerId]HostClient

	register   chan HostClient
	unregister chan NetworkPlayerId
}

func (h *NetworkHost) AddClient(id NetworkPlayerId, transport HostClientTransport) chan []byte {
	client := HostClient{}

	client.id = id
	client.transport = transport

	h.register <- client

	return h.receiver
}

func (h *NetworkHost) RemoveClient(id NetworkPlayerId) {
	h.unregister <- id
}

func (h *NetworkHost) Send(id NetworkPlayerId, data []byte) {
	h.clients[id].transport.Send(data)
}

func (h *NetworkHost) Broadcast(data []byte) {
	for i := range h.clients {
		h.clients[i].transport.Send(data)
	}
}
