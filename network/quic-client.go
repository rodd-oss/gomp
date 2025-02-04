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
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"io"
	"log"
)

type QuicClient struct {
	conn    quic.Connection
	streams []quic.Stream

	stopSignal chan empty
}

func NewQuicClient() *QuicClient {
	return &QuicClient{}
}

func (c *QuicClient) Disconnect() {
	err := c.conn.CloseWithError(0, "Connection closed")
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *QuicClient) Connect(addr string) {
	quicCfg := quic.Config{
		Tracer: qlog.DefaultConnectionTracer,
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-ebiten-ecs"},
	}

	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, &quicCfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(c quic.Connection) {
		err := c.CloseWithError(0, "Connection closed")
		if err != nil {
			log.Println(err)
		}
	}(conn)

	c.conn = conn

	// Opening Main stream
	mainStream, err := conn.OpenStreamSync(conn.Context())
	if err != nil {
		panic(err)
	}
	defer func(str quic.Stream) {
		err := str.Close()
		if err != nil {
			log.Println(err)
		}
	}(mainStream)

	c.streams = append(c.streams, mainStream)

	c.streamHandler(mainStream)
}

func (c *QuicClient) Send(msg []byte, streamId StreamId, streamIds ...StreamId) {
	streams := c.streams

	_, err := streams[streamId].Write(msg)
	if err != nil {
		log.Println(err)
	}

	for _, id := range streamIds {
		_, err := streams[id].Write(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *QuicClient) streamHandler(stream quic.Stream) {
	defer func(str quic.Stream) {
		err := str.Close()
		if err != nil {
			log.Println(err)
		}
	}(stream)

	log.Println("New stream from " + c.conn.RemoteAddr().String())
	defer log.Println("Closing stream from " + c.conn.RemoteAddr().String())

	// Reader
	for {
		msg, err := io.ReadAll(stream)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(msg))
	}
}
