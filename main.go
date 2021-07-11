package main

import (
	"webOrder/model"
	"webOrder/server"
	"webOrder/server/router"
)

func main() {

	g := router.NewEngine(server.NewServer(
		&model.Order{},
	))
	if err := g.Run(":8888"); err != nil {
		panic(err)
	}
}