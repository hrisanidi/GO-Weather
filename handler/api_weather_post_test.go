package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webServiceApp/client"
	"webServiceApp/status"
	"webServiceApp/weather"

	"github.com/stretchr/testify/mock"
)

var (
	target = "/api/weather"
)

type MockClient struct {
}

func (m *MockClient) Request(params map[string]string, url string) ([]byte, *status.Error) {
	return nil, nil
}

type MockWeather struct {
	mock.Mock
}

func (w *MockWeather) City(client client.Client, request *weather.Request) (*weather.Response, *status.Error) {
	args := w.Called(client, request)
	return args.Get(0).(*weather.Response), args.Get(1).(*status.Error)
}

func TestApiWeatherPost_InvalidField_BadRequestThrown(t *testing.T) {
	e := echo.New()
	invalidFieldJSON := `{"name":"Vlad"}`

	mockClient := new(MockClient)
	mockWeather := new(MockWeather)

	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(invalidFieldJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ApiWeatherPost(mockClient, mockWeather)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestApiWeatherPost_InvalidJSON_BadRequestThrown(t *testing.T) {
	e := echo.New()
	invalidJSON := "string"
	messageInvalidJSON := "invalid json\n"

	mockClient := new(MockClient)
	mockWeather := new(MockWeather)

	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ApiWeatherPost(mockClient, mockWeather)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, messageInvalidJSON, rec.Body.String())
	}
}

func TestApiWeatherPost_ErrorWeather(t *testing.T) {
	e := echo.New()
	validJSON := `{"city":"Brno","units":"Metric"}`
	errorMessage := "string"
	expectedError := &status.Error{Status: http.StatusUnauthorized, Error: errorMessage}
	expectedResponse := new(weather.Response)
	validRequest := &weather.Request{City: "Brno", Units: "metric"}

	mockClient := new(MockClient)
	mockWeather := new(MockWeather)

	mockWeather.On("City", mockClient, validRequest).Return(expectedResponse, expectedError)

	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(validJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ApiWeatherPost(mockClient, mockWeather)(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, errorMessage, rec.Body.String())
	}
}

func TestApiWeatherPost_Success(t *testing.T) {
	e := echo.New()
	validJSON := `{"city":"Brno","units":"Metric"}`
	var expectedError *status.Error = nil
	expectedResponse := new(weather.Response)
	validRequest := &weather.Request{City: "Brno", Units: "metric"}

	mockClient := new(MockClient)
	mockWeather := new(MockWeather)

	mockWeather.On("City", mockClient, validRequest).Return(expectedResponse, expectedError)

	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(validJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ApiWeatherPost(mockClient, mockWeather)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
