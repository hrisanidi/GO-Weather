package weather

import (
	"github.com/tidwall/gjson"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"webServiceApp/client"
	"webServiceApp/status"
)

type AwAPIFetcher interface {
	LocationToKey(client client.Client, request *Request, aw *AwWeather) (*string, *status.Error)
	KeyToConditions(client client.Client, request *Request, aw *AwWeather, locationKey string) (*Response, *status.Error)
}

type AwFetcher struct {
}

// LocationToKey get location key by city name
func (f *AwFetcher) LocationToKey(client client.Client, request *Request, aw *AwWeather) (*string, *status.Error) {
	resp, err := client.Request(map[string]string{"apikey": aw.Config.Key, "q": request.City}, aw.Config.CityURL)
	if err != nil {
		return nil, err
	}
	locationKey := gjson.GetBytes(resp, "0.Key").String()
	return &locationKey, err
}

// KeyToConditions get weather conditions by location key
func (f *AwFetcher) KeyToConditions(client client.Client, request *Request, aw *AwWeather, locationKey string) (*Response, *status.Error) {
	url := aw.Config.ConditionsURL + locationKey
	resp, err := client.Request(map[string]string{"apikey": aw.Config.Key, "details": "true"}, url)
	if err != nil {
		return nil, err
	}

	// construct the output
	formatter := cases.Title(language.English)
	out := new(Response)

	out.City = formatter.String(request.City)
	out.Temperature = gjson.GetBytes(resp, "0.Temperature."+formatter.String(request.Units)+".Value").Num
	out.Wind.Direction = gjson.GetBytes(resp, "0.Wind.Direction.English").String()
	out.Wind.Speed = gjson.GetBytes(resp, "0.Wind.Speed."+formatter.String(request.Units)+".Value").Num

	return out, nil
}
