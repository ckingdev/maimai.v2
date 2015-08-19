package main

import (
	"fmt"
	"sync"

	"github.com/cpalone/gobot"
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
	fmt.Println(b.Rooms)
	var wg sync.WaitGroup
	for _, room := range b.Rooms {
		wg.Add(1)
		go func(r *gobot.Room) {
			defer wg.Done()
			err := r.Run()
			fmt.Printf("Room %s finished running: %s", r.RoomName, err)
		}(room)
	}
	wg.Wait()
}
