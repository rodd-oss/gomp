package game

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"github.com/redis/go-redis/v9"
// )

// const (
// 	tickRate      = 100 * time.Millisecond
// 	heartbeatRate = 1 * time.Second
// 	lazyTickRate  = tickRate * 100
// )

// func loop(g *Game) {
// 	ticker := time.NewTicker(tickRate)
// 	defer ticker.Stop()

// 	heartbeat := time.NewTicker(heartbeatRate)
// 	defer heartbeat.Stop()

// 	tickerLazy := time.NewTicker(lazyTickRate)
// 	defer tickerLazy.Stop()

// 	gameStateMigration := g.db.Subscribe(context.TODO(), "game_state_migration")
// 	defer gameStateMigration.Unsubscribe(context.TODO(), "game_state_migration")

// 	for {
// 		select {
// 		case <-ticker.C:
// 			go g.mainLoop()
// 		case <-heartbeat.C:
// 			go g.heartbeatLoop()
// 		case <-tickerLazy.C:
// 			go g.lazyLoop()
// 		case msg := <-gameStateMigration.Channel():
// 			go g.syncStateLoop(msg)
// 		}
// 	}
// }

// func (g *Game) mainLoop() {
// 	if len(g.state.migrations) == 0 {
// 		return
// 	}

// 	b := make([]byte, len(g.state.migrations)*3+1)

// 	b[0] = byte(StateMigrationMessage)

// 	i := 1
// 	for _, v := range g.state.migrations {
// 		b[i] = byte(v.y)
// 		b[i+1] = byte(v.x)
// 		b[i+2] = byte(v.color)
// 		i = i + 3
// 	}

// 	err := g.db.Publish(context.TODO(), "game_state_migration", b).Err()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var mv map[string]byte = make(map[string]byte, len(g.state.migrations))
// 	for _, v := range g.state.migrations {
// 		mv[fmt.Sprintf("%d_%d", v.x, v.y)] = byte(v.color)
// 	}

// 	err = g.db.MSet(context.TODO(), mv).Err()
// 	if err != nil {
// 		fmt.Println("Error MSET:", err)
// 	}

// 	g.state.migrations = make(map[int]*StateMigration) // reset migrations

// 	// // current state as binary array
// 	// var binSatate []byte
// 	// for _, c := range g.state.current {
// 	// 	binSatate = append(binSatate, byte(c))
// 	// }

// 	// // save current state to rdb
// 	// err = g.db.Set(context.TODO(), "game_state", binSatate, 0).Err()
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// }

// func (g *Game) lazyLoop() {
// 	for _, player := range g.players {
// 		g.SendState(player)
// 	}
// }

// func (g *Game) heartbeatLoop() {
// 	for _, player := range g.players {
// 		player.SendPlayerCounter(len(g.players))
// 	}
// }

// func (g *Game) syncStateLoop(msg *redis.Message) {
// 	m := []byte(msg.Payload)

// 	for _, player := range g.players {
// 		player.SendMessage(Message{Type: websocket.BinaryMessage, Data: m})
// 	}
// }
