package main

import (
	"log"

	"github.com/omriza/go-msvc-base/pkg/adapters"
	"github.com/omriza/go-msvc-base/pkg/adapters/rest"
)

func main() {
	// Registering adapters, should happen when config loads
	rest.Register()

	// Creating driving adapters' instances
	restAdapter, err := adapters.CreateDrivingAdapter(adapters.DrivingAdapterConfig{
		"adapter-type": "rest",
	})

	if err != nil {
		log.Panicln("Aw snap...Couldn't create adapters:", err)
	}

	drivingAdapters := []adapters.DrivingAdapter{
		restAdapter,
	}

	for _, adapter := range drivingAdapters {
		log.Println("Starting adapter:", adapter.Describe())
		adapter.Start()
	}
}
