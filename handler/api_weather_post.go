package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"webServiceApp/client"
	"webServiceApp/weather"
)

func ApiWeatherPost(client client.Client, w weather.Weather) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(weather.Request)
		if err := c.Bind(request); err != nil {
			return c.String(http.StatusBadRequest, "invalid json\n")
		}

		request.Units = strings.ToLower(request.Units)
		if err := validator.New().Struct(request); err != nil {
			return c.String(http.StatusBadRequest, "invalid field value\n")
		}

		// get the weather
		output, statusError := w.City(client, request)
		if statusError != nil {
			return c.String(statusError.StatusCode(), statusError.Message())
		}

		return c.JSON(http.StatusOK, output)
	}
}
