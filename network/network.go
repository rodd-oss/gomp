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

type Mode int
type Channel int
type PeerId int

const (
	ModeNone Mode = iota
	ModeServer
	ModeClient
)

type AnyNetwork interface {
	Host(addr string)
	Connect(addr string)
	Disconnect()
	Mode() Mode
	Send(ch Channel, data []byte) error
	Receive(ch Channel) ([]byte, error)
}

type AnyClient interface {
	Disconnect()
	Send(ch Channel, data []byte) error
	Receive(ch Channel) ([]byte, error)
}

type AnyServer struct {
}
