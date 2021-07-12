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
	if err := g.Run(); err != nil {
		panic(err)
	}
}