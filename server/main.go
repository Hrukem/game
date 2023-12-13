package main

import (
	"game/game"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	world := &game.World{
		IsServer: true,
		Units:    game.Units{},
	}

	hub := newHub()
	go hub.run()

	r := gin.New()
	r.GET("/ws", func(hub *Hub, world *game.World) gin.HandlerFunc {
		return gin.HandlerFunc(func(c *gin.Context) {
			serveWs(hub, world, c.Writer, c.Request)
		})
	}(hub, world))

	if err := r.Run(":3000"); err != nil {
		log.Fatal("error Run():", err)
	}
}
