package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"sort"
	"syscall/js"
	"tomb_mates/internal/engine"

	"github.com/coder/websocket"
	"github.com/golang/protobuf/proto"
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Config struct {
	title  string
	width  int
	height int
	scale  float64
}

type Sprite struct {
	Frames []image.Image
	Frame  int
	X      float64
	Y      float64
	Side   engine.Direction
	Config image.Config
}

type Camera struct {
	X       float64
	Y       float64
	Padding float64
}

var config *Config
var world *engine.World
var camera *Camera
var frames map[string]engine.Frames
var frame int
var lastKey e.Key
var prevKey e.Key
var levelImage *e.Image

// Game implements ebiten.Game interface.
type Game struct {
	Conn *websocket.Conn
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	if world.Units[world.MyID] == nil {
		return nil
	}

	handleKeyboard(g.Conn)

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *e.Image) {
	if world.Units[world.MyID] == nil {
		return
	}

	// Write your game's rendering.
	handleCamera(screen)
	if camera == nil {
		return
	}

	frame++

	sprites := []Sprite{}
	for _, unit := range world.Units {
		sprites = append(sprites, Sprite{
			Frames: frames[unit.Skin+"_"+unit.Action].Frames,
			Frame:  int(unit.Frame),
			X:      unit.X,
			Y:      unit.Y,
			Side:   unit.Side,
			Config: frames[unit.Skin+"_"+unit.Action].Config,
		})
	}
	sort.Slice(sprites, func(i, j int) bool {
		depth1 := sprites[i].Y + float64(sprites[i].Config.Height)
		depth2 := sprites[j].Y + float64(sprites[j].Config.Height)
		return depth1 < depth2
	})

	for _, sprite := range sprites {
		op := &e.DrawImageOptions{}

		if sprite.Side == engine.Direction_left {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(sprite.Config.Width), 0)
		}

		op.GeoM.Translate(sprite.X-camera.X, sprite.Y-camera.Y)

		img := e.NewImageFromImage(sprite.Frames[(frame/7+sprite.Frame)%4])
		screen.DrawImage(img, op)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f ; FPS: %0.2f", e.ActualTPS(), e.ActualFPS()))
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

	world = &engine.World{
		Replica: true,
		Units:   map[string]*engine.Unit{},
	}

	var err error
	frames, err = engine.LoadResources()
	if err != nil {
		log.Fatal(err)
	}

	levelImage, err = prepareLevelImage()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go world.Evolve()

	url := js.Global().Get("document").Get("location").Get("origin").String()
	url = "ws" + url[4:] + "/ws"

	ws, _, err := websocket.Dial(context.TODO(), url, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func(c *websocket.Conn) {
		defer c.CloseNow()

		for {
			var _, message, err = c.Read(context.TODO())
			if err != nil {
				log.Fatal("Error reading message:", err)
			}

			event := &engine.Event{}
			world.Mx.Lock()
			err = proto.Unmarshal(message, event)
			world.Mx.Unlock()
			if err != nil {
				log.Fatal(err)
			}

			world.HandleEvent(event)

			if event.Type == engine.Event_type_connect {
				me := world.Units[world.MyID]
				camera = &Camera{
					X:       me.X,
					Y:       me.Y,
					Padding: 30,
				}
			}
		}
	}(ws)

	e.SetRunnableOnUnfocused(true)
	e.SetWindowSize(config.width, config.height)
	e.SetWindowTitle(config.title)
	game := &Game{Conn: ws}
	if err := e.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func prepareLevelImage() (*e.Image, error) {
	tileSize := 16
	level := engine.LoadLevel()
	width := len(level[0])
	height := len(level)
	levelImage := e.NewImage(width*tileSize, height*tileSize)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			op := &e.DrawImageOptions{}
			op.GeoM.Translate(float64(i*tileSize), float64(j*tileSize))

			img := e.NewImageFromImage(frames[level[j][i]].Frames[0])
			levelImage.DrawImage(img, op)
		}
	}

	return levelImage, nil
}

func handleCamera(screen *e.Image) {
	if camera == nil {
		return
	}

	player := world.Units[world.MyID]
	frame := frames[player.Skin+"_"+player.Action]
	camera.X = player.X - float64(config.width-frame.Config.Width)/2
	camera.Y = player.Y - float64(config.height-frame.Config.Height)/2

	op := &e.DrawImageOptions{}
	op.GeoM.Translate(-camera.X, -camera.Y)
	op.GeoM.Scale(1, 1)
	screen.DrawImage(levelImage, op)
}

func handleKeyboard(c *websocket.Conn) {
	event := &engine.Event{}

	if e.IsKeyPressed(e.KeyA) || e.IsKeyPressed(e.KeyLeft) {
		event = &engine.Event{
			Type: engine.Event_type_move,
			Data: &engine.Event_Move{
				Move: &engine.EventMove{
					PlayerId:  world.MyID,
					Direction: engine.Direction_left,
				},
			},
		}
		if lastKey != e.KeyA {
			lastKey = e.KeyA
		}
	}

	if e.IsKeyPressed(e.KeyD) || e.IsKeyPressed(e.KeyRight) {
		event = &engine.Event{
			Type: engine.Event_type_move,
			Data: &engine.Event_Move{
				Move: &engine.EventMove{
					PlayerId:  world.MyID,
					Direction: engine.Direction_right,
				},
			},
		}
		if lastKey != e.KeyD {
			lastKey = e.KeyD
		}
	}

	if e.IsKeyPressed(e.KeyW) || e.IsKeyPressed(e.KeyUp) {
		event = &engine.Event{
			Type: engine.Event_type_move,
			Data: &engine.Event_Move{
				Move: &engine.EventMove{
					PlayerId:  world.MyID,
					Direction: engine.Direction_up,
				},
			},
		}
		if lastKey != e.KeyW {
			lastKey = e.KeyW
		}
	}

	if e.IsKeyPressed(e.KeyS) || e.IsKeyPressed(e.KeyDown) {
		event = &engine.Event{
			Type: engine.Event_type_move,
			Data: &engine.Event_Move{
				Move: &engine.EventMove{
					PlayerId:  world.MyID,
					Direction: engine.Direction_down,
				},
			},
		}
		if lastKey != e.KeyS {
			lastKey = e.KeyS
		}
	}

	unit := world.Units[world.MyID]

	if event.Type == engine.Event_type_move {
		if prevKey != lastKey {
			world.Mx.Lock()
			message, err := proto.Marshal(event)
			world.Mx.Unlock()
			if err != nil {
				log.Println(err)
				return
			}
			err = c.Write(context.Background(), websocket.MessageBinary, message)
			if err != nil {
				// ...
				log.Fatal(err)
			}
		}
	} else {
		if unit.Action != engine.UnitActionIdle {
			event = &engine.Event{
				Type: engine.Event_type_idle,
				Data: &engine.Event_Idle{
					Idle: &engine.EventIdle{PlayerId: world.MyID},
				},
			}
			world.Mx.Lock()
			message, err := proto.Marshal(event)
			world.Mx.Unlock()
			if err != nil {
				log.Println(err)
				return
			}
			err = c.Write(context.Background(), websocket.MessageBinary, message)
			if err != nil {
				// ...
				log.Fatal(err)
			}
			lastKey = -1
		}
	}

	prevKey = lastKey
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	fmt.Println("Env not found: ", key, " - Using fallback: ", fallback)
	return fallback
}
