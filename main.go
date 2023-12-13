package main

import (
	"game/game"
	"github.com/gorilla/websocket"
	e "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"sort"
	"strconv"
)

var world game.World
var frame int
var img *e.Image
var err error

func init() {
	world = game.World{
		IsServer: false,
		Units:    game.Units{},
	}
}

func update(c *websocket.Conn) func(screen *e.Image) error {
	return func(screen *e.Image) error {
		frame++

		img, _, err = ebitenutil.NewImageFromFile(
			"sprites/frames/big_demon_idle_anim_f0.png",
			e.FilterDefault,
		)
		if err != nil {
			log.Println("error ")
		}
		if err = screen.DrawImage(img, nil); err != nil {
			return err
		}

		unitList := []*game.Unit{}
		for _, unit := range world.Units {
			unitList = append(unitList, unit)
		}
		sort.Slice(unitList, func(i, j int) bool {
			return unitList[i].Y < unitList[j].Y
		})

		for _, unit := range unitList {
			op := &e.DrawImageOptions{}
			if unit.HorizontalDirection == game.DirectionLeft {
				op.GeoM.Scale(-1, 1)
				op.GeoM.Translate(16, 0)
			}
			op.GeoM.Translate(unit.X, unit.Y)

			spriteIndex := (frame/7 + unit.Frame) % 4
			img, _, _ = ebitenutil.NewImageFromFile(
				"sprites/frames/"+
					unit.SpriteName+"_"+unit.Action+"_anim_f"+strconv.Itoa(spriteIndex)+".png",
				e.FilterDefault,
			)
			if err := screen.DrawImage(img, op); err != nil {
				log.Println("error DrawImage() in main():", err)
			}

			if e.IsKeyPressed(e.KeyD) || e.IsKeyPressed(e.KeyRight) {
				if err := c.WriteJSON(game.Event{
					Type: game.EventTypeMove,
					Data: game.EventMove{
						UnitID:    world.MyID,
						Direction: game.DirectionRight,
					},
				}); err != nil {
					log.Println("error move:", err)
				}
				return nil
			}
			if e.IsKeyPressed(e.KeyA) || e.IsKeyPressed(e.KeyLeft) {
				if err := c.WriteJSON(game.Event{
					Type: game.EventTypeMove,
					Data: game.EventMove{
						UnitID:    world.MyID,
						Direction: game.DirectionLeft,
					},
				}); err != nil {
					log.Println("error move:", err)
				}
				return nil
			}
			if e.IsKeyPressed(e.KeyW) || e.IsKeyPressed(e.KeyUp) {
				if err := c.WriteJSON(game.Event{
					Type: game.EventTypeMove,
					Data: game.EventMove{
						UnitID:    world.MyID,
						Direction: game.DirectionUp,
					},
				}); err != nil {
					log.Println("error move:", err)
				}
				return nil
			}
			if e.IsKeyPressed(e.KeyS) || e.IsKeyPressed(e.KeyDown) {
				if err := c.WriteJSON(game.Event{
					Type: game.EventTypeMove,
					Data: game.EventMove{
						UnitID:    world.MyID,
						Direction: game.DirectionDown,
					},
				}); err != nil {
					log.Println("error move:", err)
				}
				return nil
			}
			if world.Units[world.MyID].Action == game.ActionIdle {
				if err := c.WriteJSON(game.Event{
					Type: game.EventTypeIdle,
					Data: game.EventMove{
						UnitID: world.MyID,
					},
				}); err != nil {
					log.Println("error move:", err)
				}
			}
		}

		return nil
	}
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/ws", nil)
	if err != nil {
		log.Println("error Dial() in main():", err)
	}
	go func(c *websocket.Conn) {
		defer func() {
			if err := c.Close(); err != nil {
				log.Println("error close websocket:", err)
			}
		}()

		for {
			var event game.Event
			if err = c.ReadJSON(&event); err != nil {
				log.Println("error ReadJSON() in main():", err)
			}
			world.HandlerEvent(&event)

		}
	}(c)

	e.SetRunnableOnUnfocused(true)
	// Run() устарела, рекомендуется RunGame
	// настройка использования RunGame() довольно сложная
	// можно посмотреть здесь: https://github.com/hajimehoshi/ebiten/tree/main/examples
	if err := e.Run(update(c), 320, 240, 2, "game"); err != nil {
		log.Println("error e.Run():", err)
	}
}
