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
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/negrel/assert"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"io"
	"log"
	"math/big"
	"time"
)

type empty struct{}
type StreamId int

type QuicServerPeer struct {
	Id      PeerId
	conn    quic.Connection
	streams []quic.Stream
}

type QuicServer struct {
	listener   *quic.Listener
	stopSignal chan empty

	nextPeerId PeerId

	peers map[PeerId]*QuicServerPeer
}

func NewQuicServer() *QuicServer {
	return &QuicServer{
		stopSignal: make(chan empty),
		peers:      make(map[PeerId]*QuicServerPeer),
	}
}

func (s *QuicServer) Run(addr string) {
	quicCfg := quic.Config{
		Tracer: qlog.DefaultConnectionTracer,
	}

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), &quicCfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(l *quic.Listener) {
		err := l.Close()
		if err != nil {
			log.Println(err)
		}
	}(listener)

	log.Println("Listening on " + listener.Addr().String())
	defer log.Println("Stopped listening on " + listener.Addr().String())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopSignal:
			return
		default:
			conn, err := listener.Accept(ctx)
			if err != nil {
				log.Println(err)
				continue
			}

			id := s.createPeerId()
			peer := &QuicServerPeer{
				Id:      id,
				conn:    conn,
				streams: make([]quic.Stream, 0),
			}

			s.peers[id] = peer

			go s.peerHandler(peer)
		}
	}
}

func (s *QuicServer) Stop() {
	s.stopSignal <- empty{}
}

func (s *QuicServer) Broadcast(msg []byte, streamId StreamId, streamIds ...StreamId) {
	for _, peer := range s.peers {
		streams := peer.streams

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
}

func (s *QuicServer) Send(msg []byte, peerId PeerId, streamId StreamId, streamIds ...StreamId) {
	peer, ok := s.peers[peerId]
	assert.True(ok, "Peer not found")

	streams := peer.streams

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

func (s *QuicServer) peerHandler(peer *QuicServerPeer) {
	conn := peer.conn
	defer func(c quic.Connection) {
		err := c.CloseWithError(0, "Sever closed the connection")
		if err != nil {
			log.Println(err)
		}
	}(conn)

	log.Println("New connection from " + conn.RemoteAddr().String())
	defer log.Println("Closing connection from " + conn.RemoteAddr().String())

	ctxWithTimeout, cancel := context.WithTimeout(conn.Context(), 10*time.Second)
	defer cancel()

	mainStream, err := conn.AcceptStream(ctxWithTimeout)
	if err != nil {
		log.Println(err)
		return
	}
	peer.streams = append(peer.streams, mainStream)
	cancel()

	go s.streamHandler(conn, mainStream)

	for {
		newStream, err := conn.AcceptStream(conn.Context())
		if err != nil {
			log.Println(err)
			return
		}

		peer.streams = append(peer.streams, newStream)
		go s.streamHandler(conn, newStream)
	}
}

func (s *QuicServer) streamHandler(conn quic.Connection, stream quic.Stream) {
	defer func(str quic.Stream) {
		err := str.Close()
		if err != nil {
			log.Println(err)
		}
	}(stream)

	log.Println("New stream from " + conn.RemoteAddr().String())
	defer log.Println("Closing stream from " + conn.RemoteAddr().String())

	for {
		msg, err := io.ReadAll(stream)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(msg))
	}
}

func (s *QuicServer) createPeerId() PeerId {
	id := s.nextPeerId
	s.nextPeerId++
	return id
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-ebiten-ecs"},
	}
}
