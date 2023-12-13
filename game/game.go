package game

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/rand"
	"log"
	"time"
)

type Unit struct {
	ID                  string  `json:"id"`
	X                   float64 `json:"x"`
	Y                   float64 `json:"y"`
	SpriteName          string  `json:"sprite_name"`
	Action              string  `json:"action"`
	Frame               int     `json:"frame"`
	HorizontalDirection int     `json:"direction"`
}

type Units map[string]*Unit

type World struct {
	MyID     string `json:"-"`
	IsServer bool   `json:"-"`
	Units    `json:"units"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type EventConnect struct {
	Unit
}

type EventMove struct {
	UnitID    string `json:"unit_id"`
	Direction int    `json:"direction"`
}

type EventIdle struct {
	UnitID string `json:"unit_id"`
}

type EventInit struct {
	PlayerID string `json:"player_id"`
	Units    Units  `json:"units"`
}

const EventTypeConnect = "connect"
const EventTypeMove = "move"
const EventTypeIdle = "idle"
const EventTypeInit = "init"

const ActionRun = "run"
const ActionIdle = "idle"

const DirectionUp = 0
const DirectionDown = 1
const DirectionLeft = 2
const DirectionRight = 3

func (world *World) HandlerEvent(event *Event) {
	switch event.Type {
	case EventTypeConnect:
		str, err := json.Marshal(event.Data)
		if err != nil {
			log.Println("error1 marshal in World.HandlerEvent():", err)
		}
		var ev EventConnect
		if err = json.Unmarshal(str, &ev); err != nil {
			log.Println("error2 unmarshal in World.HandlerEvent():", err)
		}

		world.Units[ev.ID] = &ev.Unit

	case EventTypeInit:
		str, err := json.Marshal(event.Data)
		if err != nil {
			log.Println("error3 marshal in World.HandlerEvent():", err)
		}
		var ev EventInit
		if err = json.Unmarshal(str, &ev); err != nil {
			log.Println("error4 unmarshal in World.HandlerEvent():", err)
		}

		if !world.IsServer {
			world.MyID = ev.PlayerID
			world.Units = ev.Units
		}

	case EventTypeMove:
		str, err := json.Marshal(event.Data)
		if err != nil {
			log.Println("error5 marshal in World.HandlerEvent():", err)
		}
		var ev EventMove
		if err = json.Unmarshal(str, &ev); err != nil {
			log.Println("error6 unmarshal in World.HandlerEvent():", err)
		}

		unit := world.Units[ev.UnitID]
		unit.Action = ActionRun

		switch ev.Direction {
		case DirectionUp:
			unit.Y--
		case DirectionDown:
			unit.Y++
		case DirectionLeft:
			unit.X--
			unit.HorizontalDirection = ev.Direction
		case DirectionRight:
			unit.X++
			unit.HorizontalDirection = ev.Direction
		}

	case EventTypeIdle:
		str, err := json.Marshal(event.Data)
		if err != nil {
			log.Println("error7 marshal in World.HandlerEvent():", err)
		}
		var ev EventIdle
		if err = json.Unmarshal(str, &ev); err != nil {
			log.Println("error8 unmarshal in World.HandlerEvent():", err)
		}

		unit := world.Units[ev.UnitID]
		unit.Action = ActionIdle
	}
}

func (world *World) AddPlayer() *Unit {
	skins := []string{
		"elf_f", "elf_m", "knight_f", "knight_m",
		"lizard_f", "lizard_m", "wizzard_f", "wizzard_m",
	}
	id := uuid.NewV4().String()
	rnd := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	unit := &Unit{
		ID:         id,
		Action:     ActionIdle,
		X:          rnd.Float64() * 320,
		Y:          rnd.Float64() * 240,
		Frame:      rnd.Intn(4),
		SpriteName: skins[rnd.Intn(len(skins))],
	}

	world.Units[id] = unit

	return unit
}
