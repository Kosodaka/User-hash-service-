package main

import (
	"log"
	"mainHashService/internal/app"
	"mainHashService/internal/entity"
)

func main() {
	cfg, err := entity.NewConfig("./.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%+v\n", cfg)

	app := app.New(cfg)
	app.Run()
}
