package client

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"tomb_mates/internal/components"
	"tomb_mates/internal/game"
	"tomb_mates/internal/protos"

	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	inputs    *Inputs
	transport *Transport
	config    *Config
	Game      *game.Game

	lastUpdateTime time.Time
	dt             float64
	maxDt          float64
	avgDt          float64
	frame          int

	playerInput math.Vec2
}

func New(ctx context.Context, inputs *Inputs, transport *Transport, config *Config) *Client {
	client := &Client{
		inputs:    inputs,
		transport: transport,
		config:    config,
		Game:      game.New(true),

		lastUpdateTime: time.Now(),
		dt:             0.0,
		maxDt:          0.0,
		avgDt:          0.0,
		frame:          0,

		// TODO: Move to player controller component
		playerInput: math.Vec2{
			X: 0,
			Y: 0,
		},
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-client.Game.NetworkManager.Broadcast:
				transport.Send <- msg
			case msg := <-transport.Received:
				event := &protos.Event{}
				err := proto.Unmarshal(msg, event)
				if err != nil {
					log.Println("Error parsing message:", err)
				}

				client.Game.HandleEvent(event)
			}
		}
	}()

	return client
}

func (c *Client) Run() (err error) {
	e.SetRunnableOnUnfocused(true)
	e.SetWindowSize(c.config.Width, c.config.Height)
	e.SetWindowResizingMode(c.config.ResizingMode)
	e.SetWindowTitle(c.config.Title)

	// Connect to the server.
	err = c.transport.Connect()
	if err != nil {
		return err
	}

	if err = e.RunGame(c); err != nil {
		return err
	}

	return nil
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (c *Client) Update() error {
	g := c.Game

	if !c.transport.IsConnected {
		log.Println("No connection")
		return nil
	}

	if g.NetworkManager.MyID == nil {
		log.Println("No network id")
		return nil
	}

	if g.Entities.Units[g.NetworkManager.NetworkIdToEntityId[*g.NetworkManager.MyID]] == nil {
		log.Println("No network unit")
		return nil
	}

	c.dt = time.Now().Sub(c.lastUpdateTime).Seconds()
	if c.dt > c.maxDt {
		c.maxDt = c.dt
	}

	c.avgDt = (c.dt + c.avgDt) / 2

	c.inputs.System.Update()
	err := c.handleInput()
	if err != nil {
		log.Println(err)
		return err
	}

	g.Update(c.dt)
	c.lastUpdateTime = time.Now()

	c.frame++

	return nil
}

// TODO: Move to player controller component.
func (c *Client) handleInput() error {
	var event *protos.Event
	var pInput = c.playerInput

	input := c.inputs.Handlers[0]

	if input == nil {
		log.Println("No input handler")
		return nil
	}

	if input.ActionIsPressed(ActionMoveLeft) {
		pInput.X = -1
	} else if input.ActionIsPressed(ActionMoveRight) {
		pInput.X = 1
	} else {
		pInput.X = 0
	}

	if input.ActionIsPressed(ActionMoveUp) {
		pInput.Y = 1
	} else if input.ActionIsPressed(ActionMoveDown) {
		pInput.Y = -1
	} else {
		pInput.Y = 0
	}

	if pInput != c.playerInput {
		event = &protos.Event{
			PlayerId: *c.Game.NetworkManager.MyID,
			Type:     protos.EventType_move,
			Data: &protos.Event_Move{
				Move: &protos.EventMove{
					Direction: &protos.Vector2{
						X: pInput.X,
						Y: pInput.Y,
					},
				},
			},
		}

		c.playerInput = pInput
	}

	if event != nil {
		message, err := proto.Marshal(event)
		if err != nil {
			return err
		}

		c.transport.Send <- message
	}

	return nil
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to ad the screen size with the outside size,  return a fixed size.
func (c *Client) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (c *Client) Draw(screen *e.Image) {
	if !c.transport.IsConnected {
		return
	}

	g := c.Game

	dotGreen := e.NewImage(8, 8)
	dotGreen.Fill(color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 150,
	})

	dotRed := e.NewImage(8, 8)
	dotRed.Fill(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 150,
	})

	dotBlue := e.NewImage(32, 32)
	dotBlue.Fill(color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 150,
	})

	op := &e.DrawImageOptions{}

	// Write your game's rendering.
	// s.handleCamera(screen)
	// if camera == nil {
	// 	return
	// }

	g.Mx.Lock()
	components.Render.Each(g.World, func(e *ecs.Entry) {
		body := components.Transform.GetValue(e)

		op.GeoM.Reset()
		op.GeoM.Translate(body.LocalPosition.X, body.LocalPosition.Y)
		screen.DrawImage(dotBlue, op)
	})

	components.NetworkEntity.Each(g.World, func(e *ecs.Entry) {
		ne := components.NetworkEntity.GetValue(e)

		op.GeoM.Reset()
		op.GeoM.Translate(ne.Transform.Position.X, ne.Transform.Position.Y)
		screen.DrawImage(dotRed, op)
	})

	g.Space.EachBody(func(body *cp.Body) {
		op.GeoM.Reset()
		op.GeoM.Translate(body.Position().X, body.Position().Y)
		screen.DrawImage(dotGreen, op)
	})
	g.Mx.Unlock()

	// i := 0
	// g.Mx.Lock()
	// for _, areaEntity := range g.Entities.Areas {
	// 	println("area", areaEntity)
	// 	areaComponent := components.NetworkArea.GetValue(areaEntity)
	// 	area := areaComponent.Area

	// 	sprites[i] = &Sprite{
	// 		Frames: s.Game.Sprites[area.Skin].Frames,
	// 		Frame:  int(area.Frame),
	// 		X:      area.Position.X,
	// 		Y:      area.Position.Y,
	// 		Config: s.Game.Sprites[area.Skin].Config,
	// 	}
	// 	op.GeoM.Reset()
	// 	op.GeoM.Scale(area.Size.X/float64(sprites[i].Config.Width), area.Size.Y/float64(sprites[i].Config.Height))
	// 	sprites[i].op = op

	// 	i++
	// }
	// g.Mx.Unlock()

	// firstUnitIndex := i
	// g.Mx.Lock()
	// for _, unitEntity := range g.Entities.Units {
	// 	unitComponent := components.NetworkUnit.GetValue(unitEntity)
	// 	unit := unitComponent.Unit

	// 	sprites[i] = &Sprite{
	// 		Frames: s.Game.Sprites[unit.Skin.String()+"_"+unit.Action.String()].Frames,
	// 		Frame:  int(unit.Frame),
	// 		X:      unit.Position.X,
	// 		Y:      unit.Position.Y,
	// 		Config: s.Game.Sprites[unit.Skin.String()+"_"+unit.Action.String()].Config,
	// 		Hp:     unit.Hp,
	// 	}

	// 	op := &e.DrawImageOptions{}

	// 	if unit.Side == protos.Direction_left {
	// 		op.GeoM.Scale(-1, 1)
	// 		op.GeoM.Translate(float64(sprites[i].Config.Width), 0)
	// 	}

	// 	sprites[i].op = op

	// 	i++
	// }
	// g.Mx.Unlock()
	// hpBar := s.Game.Sprites["hp"].Frames

	// sort.Slice(sprites[firstUnitIndex:i], func(i, j int) bool {
	// 	depth1 := sprites[i].Y + float64(sprites[i].Config.Height)
	// 	depth2 := sprites[j].Y + float64(sprites[j].Config.Height)
	// 	return depth1 < depth2
	// })

	// hpOp := &e.DrawImageOptions{}
	// for _, sprite := range sprites[:i] {
	// 	if sprite.Hp > 0 {
	// 		hpOp.GeoM.Reset()
	// 		hpOp.GeoM.Scale(float64(sprite.Hp)/100.0, 1)
	// 		hpOp.GeoM.Translate(sprite.X-camera.X+float64(sprite.Config.Width)/2-16, sprite.Y-camera.Y-15)
	// 		hpFrameIndex := 4 - int(math.Ceil(float64(sprite.Hp)/25))
	// 		screen.DrawImage(hpBar[hpFrameIndex], hpOp)
	// 	}

	// 	sprite.op.GeoM.Translate(sprite.X-camera.X, sprite.Y-camera.Y)

	// 	screen.DrawImage(sprite.Frames[(frame/7+sprite.Frame)%len(sprite.Frames)], sprite.op)
	// }
	if c.config.EnableDebug {

		var debugInfo = make([]string, 0)

		debugInfo = append(debugInfo, fmt.Sprintf("TPS %0.2f", e.ActualTPS()))
		debugInfo = append(debugInfo, fmt.Sprintf("FPS %0.2f", e.ActualFPS()))
		debugInfo = append(debugInfo, fmt.Sprintf("dt %0.3f", c.dt))
		debugInfo = append(debugInfo, fmt.Sprintf("max dt %0.3f", c.maxDt))
		debugInfo = append(debugInfo, fmt.Sprintf("avg dt %0.3f", c.avgDt))
		debugInfo = append(debugInfo, fmt.Sprintf("players %d", len(g.Entities.Units)))

		if g.NetworkManager.MyID != nil {
			debugInfo = append(debugInfo, fmt.Sprintf("ID %d", *g.NetworkManager.MyID))

			myUnit := g.Entities.Units[g.NetworkManager.NetworkIdToEntityId[*g.NetworkManager.MyID]]
			if myUnit != nil {
				transform := components.Transform.GetValue(myUnit)
				debugInfo = append(debugInfo, fmt.Sprintf("posX %0.0f", transform.LocalPosition.X))
				debugInfo = append(debugInfo, fmt.Sprintf("posY %0.0f", transform.LocalPosition.Y))
			}
		}

		ebitenutil.DebugPrint(screen, strings.Join(debugInfo, "\n"))
	}
}
