package client

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
	"webServiceApp/status"
)

type Client interface {
	Request(params map[string]string, url string) ([]byte, *status.Error)
}

type RestyClient struct {
	Instance *resty.Client
}

func NewRestyClient() *RestyClient {
	client := resty.New()

	// set retry in case of 5xx errors or possible connection errors
	client.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500 || err != nil
		})

	return &RestyClient{client}
}

func (c *RestyClient) Request(params map[string]string, url string) ([]byte, *status.Error) {
	client := c.Instance
	resp, err := client.R().
		SetQueryParams(params).
		Get(url)

	if resp != nil && resp.StatusCode() >= 400 {
		return nil,
			&status.Error{Status: resp.StatusCode(), Error: "external error with code: " + resp.Status() + "\n"}
	}
	if err != nil {
		return nil,
			&status.Error{Status: http.StatusServiceUnavailable, Error: "service unavailable\n"}
	}

	return resp.Body(), nil
}
