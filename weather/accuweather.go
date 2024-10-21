package weather

import (
	"webServiceApp/client"
	"webServiceApp/config"
	"webServiceApp/status"
)

type AwWeather struct {
	Fetcher AwAPIFetcher
	Config  AwWeatherConfig
}

type AwWeatherConfig struct {
	Key           string `yaml:"key"`
	CityURL       string `yaml:"city_url"`
	ConditionsURL string `yaml:"conditions_url"`
}

func NewAwWeather() *AwWeather {
	aw := new(AwWeather)
	aw.Fetcher = new(AwFetcher)
	config.LoadConfig(&aw.Config, "accuweather.yaml")
	return aw
}

// City gets the weather for the city specified in the request
func (aw *AwWeather) City(c client.Client, r *Request) (*Response, *status.Error) {
	// get location key by city name
	locationKey, err := aw.Fetcher.LocationToKey(c, r, aw)
	if err != nil {
		return nil, err
	}

	// get weather conditions by location key
	resp := new(Response)

	resp, err = aw.Fetcher.KeyToConditions(c, r, aw, *locationKey)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
