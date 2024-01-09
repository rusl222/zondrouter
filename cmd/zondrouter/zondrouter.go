package main

import (
	"context"
	"log"

	"github.com/rusl222/zondrouter/internal/app"
	//"github.com/rusl222/zondrouter/internal/systray"
)

var Version string

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}

	<-make(chan int)

	//systray.Run()
}
