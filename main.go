package main

import (
	"github.com/carolinasolfernandez/homevision/client"
	"github.com/carolinasolfernandez/homevision/config"
	"github.com/carolinasolfernandez/homevision/house"
	"log"
)

func main() {
	config.InitConfig()

	h := house.NewHouseService(
		client.NewClient(config.EnvConfig.ClientRetries),
		config.EnvConfig.HousesUrl,
		config.EnvConfig.PhotosDir,
	)

	var doneCh = make(chan struct{})
	var errorCh = make(chan error)
	var houseCh = make(chan []house.House, config.EnvConfig.NumPages)

	go func() {
		for {
			e, ok := <-errorCh
			if ok {
				log.Printf("error: %v\n", e)
			}
		}
	}()

	go h.GetHouses(config.EnvConfig.NumPages, config.EnvConfig.NumPerPage, houseCh, errorCh)

	go h.SavePhotos(houseCh, doneCh, errorCh)

	<-doneCh
	log.Println("All houses processed")

	close(errorCh)
	close(doneCh)
}
