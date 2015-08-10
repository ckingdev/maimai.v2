package main

import (
	"github.com/cpalone/gobot/config"
	"github.com/cpalone/gobot/handlers"
	"github.com/cpalone/maimai.v2"
)


func main() {
	b, err := config.BotFromCfgFile("test.yml")
	if err != nil {
		panic(err)
	}
	for _, room := range b.Rooms {
		room.Handlers = append(room.Handlers, 
		&maimai.LinkTitleHandler{},
		 &handlers.UptimeHandler{}, 
		&handlers.PongHandler{},
		&handlers.HelpHandler{ShortDesc:"short help", LongDesc: "long help"})
	}
	b.RunAllRooms()
}
