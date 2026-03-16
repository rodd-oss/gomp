/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package tomb_mates_demo

// import (
// 	"context"
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"log"
// 	"math"
// 	"os"
// 	"strings"
// 	"syscall/js"
// 	"time"
// 	"gomp_game/internal/components"
// 	"gomp_game/internal/game"
// 	"gomp_game/internal/protos"
// 	"gomp_game/internal/resources"

// 	"github.com/coder/websocket"
// 	e "github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/jakecoffman/cp/v2"
// 	input "github.com/quasilyte/ebitengine-input"
// 	ecs "github.com/yohamta/donburi"
// 	"google.golang.org/protobuf/proto"
// )

// type Config struct {
// 	title  string
// 	width  int
// 	height int
// 	scale  float64
// }

// type Sprite struct {
// 	Frames []*e.Image
// 	Frame  int
// 	X      float64
// 	Y      float64
// 	op     *e.DrawImageOptions
// 	Config image.Config
// 	Hp     uint32
// }

// type Camera struct {
// 	X       float64
// 	Y       float64
// 	Padding float64
// }

// var camera *Camera = &Camera{
// 	X:       0,
// 	Y:       0,
// 	Padding: 30,
// }

// var frame int
// var prevKey e.Key
// var levelImage *e.Image

// // GameState implements ebiten.GameState interface.
// type GameState struct {
// 	Conn        *websocket.Conn
// 	Game        *game.Game
// 	config      *Config
// 	inputSystem input.System
// 	input       *input.Handler
// 	playerInput struct {
// 		X float64
// 		Y float64
// 	}
// }

// // Update proceeds the game state.
// // Update is called every tick (1/60 [s] by default).
// var lastUpdateTime = time.Now()
// var dt float64 = 0.0
// var maxDt float64 = 0.0
// var avgDt float64 = 0.0

// func (s *GameState) Update() error {
// 	g := s.Game

// 	if s.Conn == nil {
// 		println("No connection")

// 		return nil
// 	}

// 	if g.NetworkManager.MyID == nil {
// 		println("No network id")

// 		return nil
// 	}

// 	if g.Entities.Units[g.NetworkManager.NetworkIdToEntityId[*g.NetworkManager.MyID]] == nil {
// 		println("No network unit")
// 		return nil
// 	}

// 	dt = time.Now().Sub(lastUpdateTime).Seconds()
// 	if dt > maxDt {
// 		maxDt = dt
// 	}

// 	avgDt = (dt + avgDt) / 2

// 	s.inputSystem.Update()
// 	err := s.handleInput(s.Conn)
// 	if err != nil {
// 		println(err)
// 		return err
// 	}

// 	g.Update(dt)
// 	lastUpdateTime = time.Now()

// 	// Write your game's logical update.

// 	frame++

// 	return nil
// }

// var sprites []*Sprite

// // Draw draws the game screen.
// // Draw is called every frame (typically 1/60[s] for 60Hz display).
// func (s *GameState) Draw(screen *e.Image) {
// 	if s.Conn == nil {
// 		return
// 	}

// 	g := s.Game

// 	dotGreen := e.NewImage(8, 8)
// 	dotGreen.Fill(color.RGBA{
// 		R: 0,
// 		G: 255,
// 		B: 0,
// 		A: 150,
// 	})

// 	dotRed := e.NewImage(8, 8)
// 	dotRed.Fill(color.RGBA{
// 		R: 255,
// 		G: 0,
// 		B: 0,
// 		A: 150,
// 	})

// 	dotBlue := e.NewImage(32, 32)
// 	dotBlue.Fill(color.RGBA{
// 		R: 0,
// 		G: 0,
// 		B: 255,
// 		A: 150,
// 	})

// 	op := &e.DrawImageOptions{}

// 	// Write your game's rendering.
// 	// s.handleCamera(screen)
// 	// if camera == nil {
// 	// 	return
// 	// }

// 	g.Mx.Lock()
// 	components.RenderAssterodd.Each(g.EntityManager, func(e *ecs.Entry) {
// 		body := components.Transform.GetValue(e)

// 		op.GeoM.Reset()
// 		op.GeoM.Translate(body.LocalPosition.X, body.LocalPosition.Y)
// 		screen.DrawImage(dotBlue, op)
// 	})

// 	components.NetworkEntity.Each(g.EntityManager, func(e *ecs.Entry) {
// 		ne := components.NetworkEntity.GetValue(e)

// 		op.GeoM.Reset()
// 		op.GeoM.Translate(ne.Transform.Position.X, ne.Transform.Position.Y)
// 		screen.DrawImage(dotRed, op)
// 	})

// 	g.Space.EachBody(func(body *cp.Body) {
// 		op.GeoM.Reset()
// 		op.GeoM.Translate(body.Position().X, body.Position().Y)
// 		screen.DrawImage(dotGreen, op)
// 	})
// 	g.Mx.Unlock()

// 	// i := 0
// 	// g.Mx.Lock()
// 	// for _, areaEntity := range g.Entities.Areas {
// 	// 	println("area", areaEntity)
// 	// 	areaComponent := components.NetworkArea.GetValue(areaEntity)
// 	// 	area := areaComponent.Area

// 	// 	sprites[i] = &Sprite{
// 	// 		Frames: s.Game.Sprites[area.Skin].Frames,
// 	// 		Frame:  int(area.Frame),
// 	// 		X:      area.Position.X,
// 	// 		Y:      area.Position.Y,
// 	// 		Config: s.Game.Sprites[area.Skin].Config,
// 	// 	}
// 	// 	op.GeoM.Reset()
// 	// 	op.GeoM.Scale(area.Size.X/float64(sprites[i].Config.Width), area.Size.Y/float64(sprites[i].Config.Height))
// 	// 	sprites[i].op = op

// 	// 	i++
// 	// }
// 	// g.Mx.Unlock()

// 	// firstUnitIndex := i
// 	// g.Mx.Lock()
// 	// for _, unitEntity := range g.Entities.Units {
// 	// 	unitComponent := components.NetworkUnit.GetValue(unitEntity)
// 	// 	unit := unitComponent.Unit

// 	// 	sprites[i] = &Sprite{
// 	// 		Frames: s.Game.Sprites[unit.Skin.String()+"_"+unit.Action.String()].Frames,
// 	// 		Frame:  int(unit.Frame),
// 	// 		X:      unit.Position.X,
// 	// 		Y:      unit.Position.Y,
// 	// 		Config: s.Game.Sprites[unit.Skin.String()+"_"+unit.Action.String()].Config,
// 	// 		Hp:     unit.Hp,
// 	// 	}

// 	// 	op := &e.DrawImageOptions{}

// 	// 	if unit.Side == protos.Direction_left {
// 	// 		op.GeoM.Scale(-1, 1)
// 	// 		op.GeoM.Translate(float64(sprites[i].Config.Width), 0)
// 	// 	}

// 	// 	sprites[i].op = op

// 	// 	i++
// 	// }
// 	// g.Mx.Unlock()
// 	// hpBar := s.Game.Sprites["hp"].Frames

// 	// sort.Slice(sprites[firstUnitIndex:i], func(i, j int) bool {
// 	// 	depth1 := sprites[i].Y + float64(sprites[i].Config.Height)
// 	// 	depth2 := sprites[j].Y + float64(sprites[j].Config.Height)
// 	// 	return depth1 < depth2
// 	// })

// 	// hpOp := &e.DrawImageOptions{}
// 	// for _, sprite := range sprites[:i] {
// 	// 	if sprite.Hp > 0 {
// 	// 		hpOp.GeoM.Reset()
// 	// 		hpOp.GeoM.Scale(float64(sprite.Hp)/100.0, 1)
// 	// 		hpOp.GeoM.Translate(sprite.X-camera.X+float64(sprite.Config.Width)/2-16, sprite.Y-camera.Y-15)
// 	// 		hpFrameIndex := 4 - int(math.Ceil(float64(sprite.Hp)/25))
// 	// 		screen.DrawImage(hpBar[hpFrameIndex], hpOp)
// 	// 	}

// 	// 	sprite.op.GeoM.Translate(sprite.X-camera.X, sprite.Y-camera.Y)

// 	// 	screen.DrawImage(sprite.Frames[(frame/7+sprite.Frame)%len(sprite.Frames)], sprite.op)
// 	// }
// 	var debugInfo = make([]string, 0)

// 	debugInfo = append(debugInfo, fmt.Sprintf("TPS %0.2f", e.ActualTPS()))
// 	debugInfo = append(debugInfo, fmt.Sprintf("FPS %0.2f", e.ActualFPS()))
// 	debugInfo = append(debugInfo, fmt.Sprintf("dt %0.3f", dt))
// 	debugInfo = append(debugInfo, fmt.Sprintf("max dt %0.3f", maxDt))
// 	debugInfo = append(debugInfo, fmt.Sprintf("avg dt %0.3f", avgDt))
// 	debugInfo = append(debugInfo, fmt.Sprintf("players %d", len(g.Entities.Units)))

// 	if g.NetworkManager.MyID != nil {
// 		debugInfo = append(debugInfo, fmt.Sprintf("ID %d", *g.NetworkManager.MyID))

// 		myUnit := g.Entities.Units[g.NetworkManager.NetworkIdToEntityId[*g.NetworkManager.MyID]]
// 		if myUnit != nil {
// 			transform := components.Transform.GetValue(myUnit)
// 			debugInfo = append(debugInfo, fmt.Sprintf("posX %0.0f", transform.LocalPosition.X))
// 			debugInfo = append(debugInfo, fmt.Sprintf("posY %0.0f", transform.LocalPosition.Y))
// 		}
// 	}

// 	ebitenutil.DebugPrint(screen, strings.Join(debugInfo, "\n"))
// }

// // Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// // If you don't have to ad the screen size with the outside size,  return a fixed size.
// func (g *GameState) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
// 	return outsideWidth, outsideHeight
// }

// func main() {
// 	gameState := createState()

// 	e.SetRunnableOnUnfocused(true)
// 	e.SetWindowSize(gameState.config.width, gameState.config.height)
// 	e.SetWindowResizingMode(e.WindowResizingModeEnabled)
// 	e.SetWindowTitle(gameState.config.title)
// 	if err := e.RunGame(gameState); err != nil {
// 		log.Fatal(err)
// 	}
// }

// var Sprites = make(map[string]resources.Sprite)

// const (
// 	ActionMoveLeft input.Action = iota
// 	ActionMoveRight
// 	ActionMoveUp
// 	ActionMoveDown
// )

// func createState() (s *GameState) {
// 	var err error
// 	s = &GameState{}
// 	s.inputSystem.Init(input.SystemConfig{
// 		DevicesEnabled: input.AnyDevice,
// 	})
// 	keymap := input.Keymap{
// 		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
// 		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
// 		ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
// 		ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
// 	}
// 	s.input = s.inputSystem.NewHandler(0, keymap)

// 	s.config = &Config{
// 		title:  "Another Hero",
// 		width:  640,
// 		height: 480,
// 	}

// 	s.playerInput = struct {
// 		X float64
// 		Y float64
// 	}{0, 0}

// 	s.Game = game.New(true)

// 	Sprites, err = resources.Load()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	levelImage, err = s.prepareLevelImage()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	sprites = make([]*Sprite, s.Game.MaxPlayers+1)

// 	url := js.Global().Get("document").Get("location").Get("origin").String()
// 	url = "ws" + url[4:] + "/ws"

// 	ws, _, err := websocket.Dial(context.TODO(), url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	go func(c *websocket.Conn) {
// 		defer c.CloseNow()

// 		c.SetReadLimit(-1)

// 		for {
// 			var _, message, err = c.Read(context.TODO())
// 			if err != nil {
// 				log.Println("Error reading message:", err)
// 				return
// 			}

// 			event := &protos.Event{}
// 			err = proto.Unmarshal(message, event)
// 			if err != nil {
// 				log.Println("Error parsing message:", err)
// 				// 	continue
// 			}

// 			s.Game.HandleEvent(event)
// 		}
// 	}(ws)

// 	s.Conn = ws

// 	return s
// }

// func (s *GameState) prepareLevelImage() (*e.Image, error) {
// 	tileSize := 16
// 	level := resources.LoadLevel(25, 25)
// 	width := len(level[0])
// 	height := len(level)
// 	levelImage := e.NewImage(width*tileSize, height*tileSize)

// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			tile := level[i][j]
// 			op := &e.DrawImageOptions{}
// 			op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
// 			levelImage.DrawImage(Sprites[tile].Frames[0], op)
// 		}
// 	}

// 	return levelImage, nil
// }

// func (s *GameState) handleCamera(screen *e.Image) {
// 	if camera == nil {
// 		return
// 	}

// 	g := s.Game

// 	player := g.Entities.Units[g.NetworkManager.NetworkIdToEntityId[*g.NetworkManager.MyID]]

// 	transform := components.Transform.GetValue(player)
// 	ne := components.NetworkEntity.GetValue(player)

// 	if player == nil {
// 		return
// 	}

// 	frame := Sprites[ne.Skin.String()+"_"+"idle"]
// 	absX := camera.X - transform.LocalPosition.X + float64(s.config.width-frame.Config.Width)/2
// 	absY := camera.Y - transform.LocalPosition.Y + float64(s.config.height-frame.Config.Height)/2

// 	cameraFollowSpeed := 100000.0
// 	if math.Abs(absX) > 15 {
// 		d := (absX * absX * absX / cameraFollowSpeed)
// 		d = roundFloat(d, 1)
// 		camera.X = camera.X - d
// 	}
// 	if math.Abs(absY) > 15 {
// 		d := (absY * absY * absY / cameraFollowSpeed)
// 		d = roundFloat(d, 1)
// 		camera.Y = camera.Y - d
// 	}

// 	op := &e.DrawImageOptions{}
// 	op.GeoM.Translate(-camera.X, -camera.Y)
// 	op.GeoM.Scale(1, 1)

// 	screen.DrawImage(levelImage, op)
// }

// func roundFloat(f float64, n int) float64 {
// 	m := math.Pow(10, float64(n))
// 	return math.Round(f*m) / m
// }

// func (s *GameState) handleInput(c *websocket.Conn) error {
// 	g := s.Game

// 	var event *protos.Event
// 	pInput := s.playerInput

// 	if s.input.ActionIsPressed(ActionMoveLeft) {
// 		pInput.X = -1
// 	} else if s.input.ActionIsPressed(ActionMoveRight) {
// 		pInput.X = 1
// 	} else {
// 		pInput.X = 0
// 	}

// 	if s.input.ActionIsPressed(ActionMoveUp) {
// 		pInput.Y = 1
// 	} else if s.input.ActionIsPressed(ActionMoveDown) {
// 		pInput.Y = -1
// 	} else {
// 		pInput.Y = 0
// 	}

// 	if pInput != s.playerInput {
// 		event = &protos.Event{
// 			Type: protos.EventType_move,
// 			Data: &protos.Event_Move{
// 				Move: &protos.EventMove{
// 					Direction: &protos.Vector2{
// 						X: pInput.X,
// 						Y: pInput.Y,
// 					},
// 				},
// 			},
// 			PlayerId: *g.NetworkManager.MyID,
// 		}

// 		s.playerInput = pInput
// 	}

// 	if event != nil {
// 		g.HandleEvent(event)

// 		message, err := proto.Marshal(event)
// 		if err != nil {
// 			return err
// 		}

// 		err = c.Write(context.Background(), websocket.MessageBinary, message)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func getEnv(key, fallback string) string {
// 	if value, ok := os.LookupEnv(key); ok {
// 		return value
// 	}

// 	log.Println("Env not found: ", key, " - Using fallback: ", fallback)
// 	return fallback
// }
