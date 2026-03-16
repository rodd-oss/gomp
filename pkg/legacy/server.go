/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package legacy

//import (
//	"fmt"
//	"os"
//	"runtime"
//	"strings"
//	"sync"
//	"time"
//
//	"net/http"
//	_ "net/http/pprof"
//
//	"github.com/gorilla/sessions"
//	"github.com/gorilla/websocket"
//	"github.com/labstack/echo-contrib/session"
//	"github.com/labstack/echo/v4"
//	"github.com/labstack/echo/v4/middleware"
//	"github.com/labstack/gommon/log"
//	echopprof "github.com/sevenNt/echo-pprof"
//	"golang.org/x/time/rate"
//)
//
//func (game *Game) RunServer() {
//	runtime.SetCPUProfileRate(1000)
//
//	e := echo.New()
//	// e.Use(middleware.Logger())
//	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
//		StackSize:         1 << 10, // 1 KB
//		LogLevel:          log.ERROR,
//		DisableStackAll:   true,
//		DisablePrintStack: true,
//	}))
//
//	e.Use(middleware.BodyLimit("2M"))
//	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(60))))
//	e.Use(session.Middleware(sessions.NewCookieStore([]byte(getEnv("AUTH_SECRET", "jdkljskldjslk")))))
//	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
//		Level: 5,
//		Skipper: func(c echo.Context) bool {
//			return strings.Contains(c.Path(), "ws") // Change "metrics" for your own path
//		},
//	}))
//
//	e.GET("/ws", wsHandler(game))
//
//	echopprof.Wrap(e)
//
//	go e.Start(":27015")
//}
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//func wsHandler(game *Game) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		if len(game.Network.Host.clients) >= game.Network.Host.MaxClients {
//			return c.String(http.StatusBadRequest, "Too many players")
//		}
//
//		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
//		if err != nil {
//			return err
//		}
//
//		return handleWsConnection(c, conn, game)
//	}
//}
//
//type wsClient struct {
//	id       NetworkPlayerId
//	conn     *websocket.Conn
//	send     chan []byte
//	receiver chan []byte
//}
//
//func (c *wsClient) Send(msg []byte) {
//	c.send <- msg
//}
//
//func (c *wsClient) Close() error {
//	return c.conn.Close()
//}
//
//func handleWsConnection(ctx echo.Context, conn *websocket.Conn, game *Game) error {
//	var id NetworkPlayerId = 1
//
//	client := wsClient{}
//
//	client.id = id
//	client.conn = conn
//	client.send = make(chan []byte, 512)
//
//	receiver := game.Network.Host.AddClient(id, &client)
//	defer game.Network.Host.RemoveClient(id)
//
//	client.receiver = receiver
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	go client.readPump(&wg)
//	go client.writePump(&wg)
//
//	wg.Wait()
//
//	return conn.Close()
//}
//
//const (
//	// Time allowed to write a message to the peer.
//	writeWait = 10 * time.Second
//
//	// Time allowed to read the next pong message from the peer.
//	pongWait = 60 * time.Second
//
//	// Send pings to peer with this period. Must be less than pongWait.
//	pingPeriod = (pongWait * 9) / 10
//
//	// Maximum message size allowed from peer.
//	maxMessageSize = 512
//)
//
//// readPump pumps messages from the websocket connection to the hub.
////
//// The application runs readPump in a per-connection goroutine. The application
//// ensures that there is at most one reader on a connection by executing all
//// reads from this goroutine.
//func (c *wsClient) readPump(wg *sync.WaitGroup) {
//	defer func() {
//		wg.Done()
//	}()
//
//	c.conn.SetReadLimit(maxMessageSize)
//	c.conn.SetReadDeadline(time.Now().Add(pongWait))
//	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
//
//	for {
//		_, message, err := c.conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				fmt.Printf("error: %v", err)
//			}
//			return
//		}
//
//		c.receiver <- message
//	}
//}
//
//// writePump pumps messages from the hub to the websocket connection.
////
//// A goroutine running writePump is started for each connection. The
//// application ensures that there is at most one writer to a connection by
//// executing all writes from this goroutine.
//func (c *wsClient) writePump(wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	pingTicker := time.NewTicker(pingPeriod)
//	defer pingTicker.Stop()
//
//	for {
//		select {
//		case message, ok := <-c.send:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if !ok {
//				// The hub closed the channel.
//				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			err := c.conn.WriteMessage(websocket.BinaryMessage, message)
//			if err != nil {
//				return
//			}
//
//		case <-pingTicker.C:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				return
//			}
//		}
//	}
//}
//
//func getEnv(key, fallback string) string {
//	if value, ok := os.LookupEnv(key); ok {
//		return value
//	}
//	return fallback
//}
