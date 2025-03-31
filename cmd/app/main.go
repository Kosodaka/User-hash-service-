package main

import (
	"log"
	"mainHashService/app/entity"
	"mainHashService/cmd"
)

func main() {
	cfg, err := entity.NewConfig("./.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%+v\n", cfg)

	app := cmd.New(cfg)
	app.Run()
}
