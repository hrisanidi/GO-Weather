package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	params  = new(map[string]string)
	mockURL = "https://example.com"
)

func TestRestyClient_StatusOkThrown(t *testing.T) {
	rst := resty.New()
	httpmock.ActivateNonDefault(rst.GetClient())
	defer httpmock.DeactivateAndReset()
	mockClient := RestyClient{Instance: rst}
	success := "success"

	httpmock.RegisterResponder("GET", mockURL,
		httpmock.NewStringResponder(http.StatusOK, success))
	resp1, err1 := mockClient.Request(*params, mockURL)
	assert.Equal(t, success, string(resp1))
	assert.Nil(t, err1)
}

func TestRestyClient_BadRequestThrown(t *testing.T) {
	rst := resty.New()
	httpmock.ActivateNonDefault(rst.GetClient())
	defer httpmock.DeactivateAndReset()
	mockClient := RestyClient{Instance: rst}

	httpmock.RegisterResponder("GET", mockURL,
		httpmock.NewStringResponder(http.StatusBadRequest, ""))
	resp2, err2 := mockClient.Request(*params, "https://example.com")
	assert.Nil(t, resp2)
	assert.Equal(t, http.StatusBadRequest, err2.StatusCode())
}

func TestRestyClient_InternalError(t *testing.T) {
	rst := resty.New()
	httpmock.ActivateNonDefault(rst.GetClient())
	defer httpmock.DeactivateAndReset()
	mockClient := RestyClient{Instance: rst}

	httpmock.RegisterResponder("POST", "https://example.com/test",
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("")
		},
	)
	resp3, err3 := mockClient.Request(*params, "https://example.com/test")
	assert.Nil(t, resp3)
	assert.Equal(t, http.StatusServiceUnavailable, err3.StatusCode())
}
