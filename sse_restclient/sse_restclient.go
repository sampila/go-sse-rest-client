package sse_restclient

import (
	"strings"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
	"encoding/json"
	"github.com/sampila/go-utils/rest_errors"
	"errors"
)

/*const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-Caller-Id"

	paramAccessToken = "access_token"
)*/

var (
	sseRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:9011",
		Timeout: 3000 * time.Millisecond,
	}
)

type EventRequest struct {
  Data map[string]string `json:"data"`
}

func (p *EventRequest) SentEvent()(interface{}, rest_errors.RestErr) {
	response := sseRestClient.Post("/say", p)
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to upload file",
			errors.New("network timeout"))
	}
	if response.StatusCode > 299 {
		if strings.TrimSpace(response.String()) == "expired access token" {
			restErr := rest_errors.NewUnauthorizedError(response.String())
			return nil, restErr
		}
		restErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to upload file", err)
		}
		return nil, restErr
	}

	var respBody interface{}
	if err := json.Unmarshal(response.Bytes(), &respBody); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal response",
			errors.New("error processing json"))
	}
	return &respBody, nil
}
