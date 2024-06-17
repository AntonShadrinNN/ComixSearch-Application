// Package main is the entry point of the program, where execution begins.
package main

import (
	_ "comixsearch/docs"
	"comixsearch/internal/app"
	"comixsearch/internal/config"
	"comixsearch/internal/ports/httpgin"
	"context"
	"log"
)

//	@title			Comic search application documentation
//	@version		1.0.0
//	@description	A collection of endpoints available to retrieve the comices for a specific keywords.

//	@contact.name	Maintainers
//	@contact.url	https://github.com/AntonShadrinNN/ComixSearch-Application.git
//	@contact.email	svebo3348@gmail.com
//	@host			localhost:8080
//	@accept			json
//	@produce		json
//	@schemes		http

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
