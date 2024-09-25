package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"sort"
	"syscall/js"
	"time"
	"tomb_mates/internal/game"
	"tomb_mates/internal/protos"
	"tomb_mates/internal/resources"

	"github.com/coder/websocket"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	title  string
	width  int
	height int
	scale  float64
}

type Sprite struct {
	Frames []*e.Image
	Frame  int
	X      float64
	Y      float64
	Side   protos.Direction
	Config image.Config
}

type Camera struct {
	X       float64
	Y       float64
	Padding float64
}

var config *Config
var world *game.Game
var camera *Camera = &Camera{
	X:       0,
	Y:       0,
	Padding: 30,
}
var frames map[string]resources.Frames
var frame int
var currentKey e.Key
var prevKey e.Key
var levelImage *e.Image

// Game implements ebiten.Game interface.
type Game struct {
	Conn *websocket.Conn
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
var lastUpdateTime = time.Now()
var dt float64 = 0.0
var maxDt float64 = 0.0
var avgDt float64 = 0.0

var traficIn = 0.0

func (g *Game) Update() error {
	dt = time.Now().Sub(lastUpdateTime).Seconds()
	if dt > maxDt {
		maxDt = dt
	}

	avgDt = (dt + avgDt) / 2

	world.HandlePhysics(dt)
	lastUpdateTime = time.Now()

	// Write your game's logical update.
	if world.Units[world.MyID] == nil {
		return nil
	}

	err := handleInput(g.Conn)
	if err != nil {
		return err
	}

	return nil
}

var sprites = make([]*Sprite, 1024)

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *e.Image) {
	l := len(world.Units)
	if l == 0 {
		log.Println("No units")
		return
	}

	// Write your game's rendering.
	handleCamera(screen)
	if camera == nil {
		return
	}

	frame++

	i := 0
	world.Mx.Lock()
	for _, unit := range world.Units {
		sprites[i] = &Sprite{
			Frames: frames[unit.Skin.String()+"_"+unit.Action.String()].Frames,
			Frame:  int(unit.Frame),
			X:      unit.Position.X,
			Y:      unit.Position.Y,
			Side:   unit.Side,
			Config: frames[unit.Skin.String()+"_"+unit.Action.String()].Config,
		}
		i++
	}
	world.Mx.Unlock()

	sort.Slice(sprites[:i], func(i, j int) bool {
		// if sprites[i] == nil || sprites[j] == nil {
		// 	return true
		// }
		depth1 := sprites[i].Y + float64(sprites[i].Config.Height)
		depth2 := sprites[j].Y + float64(sprites[j].Config.Height)
		return depth1 < depth2
	})

	for _, sprite := range sprites[:i] {
		op := &e.DrawImageOptions{}

		if sprite.Side == protos.Direction_left {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(sprite.Config.Width), 0)
		}

		op.GeoM.Translate(sprite.X-camera.X, sprite.Y-camera.Y)

		screen.DrawImage(sprite.Frames[(frame/7+sprite.Frame)%4], op)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f ; FPS: %0.2f ; dt: %0.3f ; maxdt: %0.3f ; avgdt: %0.3f ; players: %d ; posX: %0.3f ; posY: %0.3f", e.ActualTPS(), e.ActualFPS(), dt, maxDt, avgDt, len(world.Units), world.Units[world.MyID].Position.X, world.Units[world.MyID].Position.Y))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func init() {
	config = &Config{
		title:  "Another Hero",
		width:  640,
		height: 480,
	}

	var err error
	frames, err = resources.Load()
	if err != nil {
		log.Fatal(err)
	}

	levelImage, err = prepareLevelImage()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	world = game.New(true, map[uint32]*protos.Unit{})

	url := js.Global().Get("document").Get("location").Get("origin").String()
	url = "ws" + url[4:] + "/ws"

	ws, _, err := websocket.Dial(context.TODO(), url, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func(c *websocket.Conn) {
		defer c.CloseNow()

		c.SetReadLimit(-1)

		for {
			var _, message, err = c.Read(context.TODO())
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			event := &protos.Event{}
			err = proto.Unmarshal(message, event)
			if err != nil {
				log.Println("Error parsing message:", err)
				continue
			}

			world.HandleEvent(event)
		}
	}(ws)

	e.SetRunnableOnUnfocused(true)
	e.SetWindowSize(config.width, config.height)
	e.SetWindowResizingMode(e.WindowResizingModeEnabled)
	e.SetWindowTitle(config.title)
	game := &Game{Conn: ws}
	if err := e.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func prepareLevelImage() (*e.Image, error) {
	tileSize := 16
	level := resources.LoadLevel()
	width := len(level[0])
	height := len(level)
	levelImage := e.NewImage(width*tileSize, height*tileSize)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			op := &e.DrawImageOptions{}
			op.GeoM.Translate(float64(i*tileSize), float64(j*tileSize))

			levelImage.DrawImage(frames[level[j][i]].Frames[0], op)
		}
	}

	return levelImage, nil
}

func handleCamera(screen *e.Image) {
	if camera == nil {
		return
	}

	player := world.Units[world.MyID]
	if player == nil {
		return
	}

	frame := frames[player.Skin.String()+"_"+player.Action.String()]
	camera.X = player.Position.X - float64(config.width-frame.Config.Width)/2
	camera.Y = player.Position.Y - float64(config.height-frame.Config.Height)/2

	op := &e.DrawImageOptions{}
	op.GeoM.Translate(-camera.X, -camera.Y)
	op.GeoM.Scale(1, 1)

	screen.DrawImage(levelImage, op)
}

func handleInput(c *websocket.Conn) error {
	event := &protos.Event{}

	if e.IsKeyPressed(e.KeyA) || e.IsKeyPressed(e.KeyLeft) {
		event = &protos.Event{
			Type: protos.EventType_move,
			Data: &protos.Event_Move{
				Move: &protos.EventMove{
					Direction: protos.Direction_left,
				},
			},
			PlayerId: world.MyID,
		}
		if currentKey != e.KeyA {
			currentKey = e.KeyA
		}
	}

	if e.IsKeyPressed(e.KeyD) || e.IsKeyPressed(e.KeyRight) {
		event = &protos.Event{
			Type: protos.EventType_move,
			Data: &protos.Event_Move{
				Move: &protos.EventMove{
					Direction: protos.Direction_right,
				},
			},
			PlayerId: world.MyID,
		}
		if currentKey != e.KeyD {
			currentKey = e.KeyD
		}
	}

	if e.IsKeyPressed(e.KeyW) || e.IsKeyPressed(e.KeyUp) {
		event = &protos.Event{
			Type: protos.EventType_move,
			Data: &protos.Event_Move{
				Move: &protos.EventMove{
					Direction: protos.Direction_up,
				},
			},
			PlayerId: world.MyID,
		}
		if currentKey != e.KeyW {
			currentKey = e.KeyW
		}
	}

	if e.IsKeyPressed(e.KeyS) || e.IsKeyPressed(e.KeyDown) {
		event = &protos.Event{
			Type: protos.EventType_move,
			Data: &protos.Event_Move{
				Move: &protos.EventMove{
					Direction: protos.Direction_down,
				},
			},
			PlayerId: world.MyID,
		}
		if currentKey != e.KeyS {
			currentKey = e.KeyS
		}
	}

	unit := world.Units[world.MyID]

	if event.Type == protos.EventType_move {
		if prevKey != currentKey {
			message, err := proto.Marshal(event)
			if err != nil {
				return err
			}

			world.HandleEvent(event)

			err = c.Write(context.Background(), websocket.MessageBinary, message)
			if err != nil {
				return err
			}
		}
	} else {
		if unit.Action != protos.Action_idle {
			event = &protos.Event{
				Type:     protos.EventType_stop,
				PlayerId: world.MyID,
			}

			world.HandleEvent(event)

			message, err := proto.Marshal(event)
			if err != nil {
				return err
			}
			err = c.Write(context.Background(), websocket.MessageBinary, message)
			if err != nil {
				// ...
				return err
			}
			currentKey = -1
		}
	}

	prevKey = currentKey

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	fmt.Println("Env not found: ", key, " - Using fallback: ", fallback)
	return fallback
}
