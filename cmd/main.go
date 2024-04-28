// Package main is the entry point of the program, where execution begins.
package main

import (
	"comixsearch/internal/app"
	"comixsearch/internal/config"
	"comixsearch/internal/ports/httpgin"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	// get app configuration
	confPath := "config.yml"
	conf, err := config.GetConfig(confPath)
	if err != nil {
		log.Fatalf("Enable to get configuration. Error: %s", err)
		return
	}

	a, err := app.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("Enable to create app. Error %s", err)
	}

	err = httpgin.Run(ctx, a, conf.Port)
	if err != nil {
		log.Printf("Run server error %s", err)
		return
	}

}
