package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	quicCfg := quic.Config{
		Tracer: qlog.DefaultConnectionTracer,
	}
	conn, err := quic.DialAddr(context.Background(), "127.0.0.1:1234", tlsConf, &quicCfg)
	if err != nil {
		panic(err)
	}
	defer conn.CloseWithError(0, "defer close")

	go handleConn(conn)
	handleConn(conn)
	for {

	}
}

func handleConn(conn quic.Connection) {
	str, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}
	defer str.Close()

	_, err = str.Write([]byte("Hello, world!"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent")
	for {
	}
}
