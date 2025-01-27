package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"log"
	"math/big"
	"time"
)

func main() {
	quicCfg := quic.Config{
		Tracer: qlog.DefaultConnectionTracer,
	}
	listener, err := quic.ListenAddr("127.0.0.1:1234", generateTLSConfig(), &quicCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	ctx := context.Background()

	log.Println("Listening on " + listener.Addr().String())
	for {
		conn, err := listener.Accept(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Println("Accepted connection from " + conn.RemoteAddr().String())
		go handleConnection(conn)
		// handle the connection, usually in a new Go routine
	}
}

func handleConnection(conn quic.Connection) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		str, err := conn.AcceptStream(ctxWithTimeout)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("Accepted stream from " + conn.RemoteAddr().String())

		go handleStream(str)
	}
}

func handleStream(str quic.Stream) {
	defer str.Close()

	for {
		buf := make([]byte, 1024)
		n, err := str.Read(buf)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("Received " + string(buf[:n]))
	}
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
		NextProtos:   []string{"quic-echo-example"},
	}
}
