/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package network

import (
	"github.com/negrel/assert"
	"github.com/quic-go/quic-go"
)

var Quic = &QuicNetwork{}

type QuicNetwork struct {
	mode Mode

	// Server-side
	server *QuicServer
	client *QuicClient
}

// Host is Server-side method to host the server
func (n *QuicNetwork) Host(addr string) {
	assert.True(n.mode == ModeNone, "QuicNetwork is already in server mode")

	n.server = NewQuicServer()
	go n.server.Run(addr)
	n.mode = ModeServer
}

// Stop is Server-side method to Stop the server
func (n *QuicNetwork) Stop() {
	assert.True(n.mode == ModeServer, "QuicNetwork is not in server mode")

	n.server.Stop()
	n.server = nil
	n.mode = ModeNone
}

// Connect is Client-side method to connect to the server
func (n *QuicNetwork) Connect(addr string) {
	assert.True(n.mode == ModeNone, "QuicNetwork is already in use.")

	n.client = NewQuicClient()
	go n.client.Connect(addr)
	n.mode = ModeClient
}

// Disconnect is Client-side method to disconnect from the server
func (n *QuicNetwork) Disconnect() {
	assert.True(n.mode == ModeClient, "QuicNetwork is not in client mode")

	n.client.Disconnect()
	n.client = nil
	n.mode = ModeNone
}

func (n *QuicNetwork) Mode() Mode {
	return n.mode
}

// Send is both Server-side and Client-side method to send data to all peers
func (n *QuicNetwork) Send(data []byte, streamId StreamId, streamIds ...StreamId) error {
	assert.True(n.mode != ModeNone, "QuicNetwork is not in use")

	switch n.mode {
	case ModeNone:
		return quic.ErrTransportClosed
	case ModeServer:
		n.server.Broadcast(data, streamId, streamIds...)
		return nil
	case ModeClient:
		n.client.Send(data, streamId, streamIds...)
		return nil
	default:
		panic("QuicNetwork is in unknown mode")
	}
}

// SendTo is Server-side method to send data to a specific peer
func (n *QuicNetwork) SendTo(peer PeerId, data []byte, streamId StreamId, streamIds ...StreamId) error {
	assert.True(n.mode == ModeServer, "QuicNetwork is not in server mode")

	n.server.Send(data, peer, streamId, streamIds...)

	return nil
}

// Receive is both Server-side and Client-side method to receive data from all peers
func (n *QuicNetwork) Receive(ch Channel) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
