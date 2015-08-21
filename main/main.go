package main

import (
	"github.com/cpalone/gobot/config"
	"github.com/cpalone/maimai.v2"
)

func main() {
	b, err := config.BotFromCfgFile("MaiMai.yml")
	if err != nil {
		panic(err)
	}
	for _, room := range b.Rooms {
		room.Handlers = append(room.Handlers,
			&maimai.LinkTitleHandler{},
		)
	}
	b.RunAllRooms()
}
