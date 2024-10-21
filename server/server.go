package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"time"
	"webServiceApp/client"
	"webServiceApp/config"
	"webServiceApp/handler"
	"webServiceApp/weather"
)

type Config struct {
	Port string `yaml:"port"`
}

func NewSeverConfig() *Config {
	c := &Config{}
	config.LoadConfig(c, "server.yaml")
	return c
}

func NewServer() *echo.Echo {
	e := echo.New()
	return e
}

func AddServerRoutes(e *echo.Echo, client client.Client, weather weather.Weather) {
	e.POST("/api/weather", handler.ApiWeatherPost(client, weather))
}

// StartServer starts the server and sets up graceful shutdown
func StartServer(e *echo.Echo, config *Config) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(config.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
