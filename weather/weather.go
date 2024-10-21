package weather

import (
	"webServiceApp/client"
	"webServiceApp/status"
)

type Weather interface {
	City(client client.Client, request *Request) (*Response, *status.Error)
}

type Request struct {
	City  string `json:"city" validate:"required"`
	Units string `json:"units" validate:"required,eq=metric|eq=imperial"`
}

type Response struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Wind        Wind    `json:"wind"`
}

type Wind struct {
	Direction string  `json:"direction"`
	Speed     float64 `json:"speed"`
}
