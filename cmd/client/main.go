package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"sort"
	"strings"
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
	Hp     uint32
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

var sprites []*Sprite

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
			Hp:     unit.Hp,
		}
		i++
	}
	world.Mx.Unlock()
	hpBar := frames["hp"].Frames

	sort.Slice(sprites[:i], func(i, j int) bool {
		depth1 := sprites[i].Y + float64(sprites[i].Config.Height)
		depth2 := sprites[j].Y + float64(sprites[j].Config.Height)
		return depth1 < depth2
	})
	hpOp := &e.DrawImageOptions{}

	for _, sprite := range sprites[:i] {
		if sprite.Hp > 0 {
			hpOp.GeoM.Reset()
			hpOp.GeoM.Scale(float64(sprite.Hp)/100.0, 1)
			hpOp.GeoM.Translate(sprite.X-camera.X+float64(sprite.Config.Width)/2-16, sprite.Y-camera.Y-15)
			hpFrameIndex := 4 - int(math.Ceil(float64(sprite.Hp)/25))
			fmt.Println(hpFrameIndex)
			screen.DrawImage(hpBar[hpFrameIndex], hpOp)
		}

		op := &e.DrawImageOptions{}

		if sprite.Side == protos.Direction_left {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(sprite.Config.Width), 0)
		}

		op.GeoM.Translate(sprite.X-camera.X, sprite.Y-camera.Y)

		screen.DrawImage(sprite.Frames[(frame/7+sprite.Frame)%4], op)
	}
	var debugInfo = make([]string, 0)

	debugInfo = append(debugInfo, fmt.Sprintf("TPS %0.2f", e.ActualTPS()))
	debugInfo = append(debugInfo, fmt.Sprintf("FPS %0.2f", e.ActualFPS()))
	debugInfo = append(debugInfo, fmt.Sprintf("dt %0.3f", dt))
	debugInfo = append(debugInfo, fmt.Sprintf("max dt %0.3f", maxDt))
	debugInfo = append(debugInfo, fmt.Sprintf("avg dt %0.3f", avgDt))
	debugInfo = append(debugInfo, fmt.Sprintf("players %d", len(world.Units)))

	myUnit := world.Units[world.MyID]
	if myUnit != nil {
		debugInfo = append(debugInfo, fmt.Sprintf("posX %0.0f", myUnit.Position.X))
		debugInfo = append(debugInfo, fmt.Sprintf("posY %0.0f", myUnit.Position.Y))
	}

	ebitenutil.DebugPrint(screen, strings.Join(debugInfo, "\n"))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func init() {

}

func main() {
	var err error

	config = &Config{
		title:  "Another Hero",
		width:  640,
		height: 480,
	}

	frames, err = resources.Load()
	if err != nil {
		log.Fatal(err)
	}

	levelImage, err = prepareLevelImage()
	if err != nil {
		log.Fatal(err)
	}

	world = game.New(true, map[uint32]*protos.Unit{})
	sprites = make([]*Sprite, world.MaxPlayers+1)

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
	level := resources.LoadLevel(25, 25)
	width := len(level[0])
	height := len(level)
	levelImage := e.NewImage(width*tileSize, height*tileSize)

	log.Println(width, height)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			tile := level[i][j]
			op := &e.DrawImageOptions{}
			op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
			levelImage.DrawImage(frames[tile].Frames[0], op)
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
	absX := camera.X - player.Position.X + float64(config.width-frame.Config.Width)/2
	absY := camera.Y - player.Position.Y + float64(config.height-frame.Config.Height)/2

	cameraFollowSpeed := 100000.0
	if math.Abs(absX) > 15 {
		d := (absX * absX * absX / cameraFollowSpeed)
		d = roundFloat(d, 1)
		camera.X = camera.X - d
	}
	if math.Abs(absY) > 15 {
		d := (absY * absY * absY / cameraFollowSpeed)
		d = roundFloat(d, 1)
		camera.Y = camera.Y - d
	}

	op := &e.DrawImageOptions{}
	op.GeoM.Translate(-camera.X, -camera.Y)
	op.GeoM.Scale(1, 1)

	screen.DrawImage(levelImage, op)
}

func roundFloat(f float64, n int) float64 {
	m := math.Pow(10, float64(n))
	return math.Round(f*m) / m
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
