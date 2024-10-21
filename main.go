package main

import (
	"github.com/defval/di"
	"log"
	"webServiceApp/client"
	"webServiceApp/server"
	"webServiceApp/weather"
)

func main() {
	di.SetTracer(&di.StdTracer{})

	container, err := di.New(
		di.Provide(client.NewRestyClient, di.As(new(client.Client))),
		di.Provide(weather.NewAwWeather, di.As(new(weather.Weather))),
		di.Provide(server.NewSeverConfig),
		di.Provide(server.NewServer),
		di.Invoke(server.AddServerRoutes),
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err = container.Invoke(server.StartServer); err != nil {
		log.Fatalln(err)
	}
}
