package hostingde

// Adapted from hostingde provider in https://github.com/go-acme/lego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://secure.hosting.de/api/dns/v1/json"

// Client -
type Client struct {
	HTTPClient *http.Client
	accountId  string
	authToken  string
	baseURL    string
}

func NewClient(accountId, authToken *string) *Client {
	var account, token string

	if accountId != nil {
		account = *accountId
	}
	if authToken != nil {
		token = *authToken
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		accountId:  account,
		authToken:  token,
		baseURL:    defaultBaseURL,
	}

	return &c
}

func (c *Client) doRequestIter(httpMethod string, uri string, request Request, response interface{}, iteration int) ([]byte, error) {
	if iteration > 8 {
		return nil, fmt.Errorf("reached max retry count, status of ZoneConfig in response is still blocked")
	}
	if request.getAuthToken() == "" {
		request.setAuthToken(c.authToken)
	}
	if request.getAccountId() == "" {
		request.setAccountId(c.accountId)
	}

	rawBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(rawBody))
	req, err := http.NewRequest(httpMethod, uri, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying API: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(toErrorWithNewlines(uri, body))
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, toErrorWithNewlines(uri, body))
	}

	// Sometimes the API returns an undocumented blocked status
	var br *BaseResponse
	switch r := response.(type) {
	default:
		return nil, fmt.Errorf("error: invalid response type: %T", r)
	case *ZoneUpdateResponse:
		br = &r.BaseResponse
	case *ZoneCreateResponse:
		br = &r.BaseResponse
	case *ZoneDeleteResponse:
		br = &r.BaseResponse
	case *ZoneConfigsFindResponse:
		if len(r.Response.Data) == 0 {
			return nil, fmt.Errorf("%v: uri: %s %s", err, uri, response)
		}
		br = &r.BaseResponse
	case *ZonesFindResponse:
		if len(r.Response.Data) == 0 {
			return nil, fmt.Errorf("%v: uri: %s %s", err, uri, response)
		}
		br = &r.BaseResponse
	case *RecordsFindResponse:
		br = &r.BaseResponse
	case *RecordsUpdateResponse:
		br = &r.BaseResponse
	}

	iteration++

	// The API returns two status strings:
	// https://www.hosting.de/api/#responses
	// https://www.hosting.de/api/#the-zoneconfig-object
	// If the first is an error, we check if it's because the resource is blocked
	if br.Status == "error" {
		var blocked bool
		for _, err := range br.Errors {
			if err.Value == "blocked" {
				blocked = true
			}
		}
		if blocked {
			fmt.Printf("Request blocked, triggering new request: %d", iteration)
			time.Sleep(1 * time.Second)
			return c.doRequestIter(httpMethod, uri, request, response, iteration)
		}
	}

	return body, err
}

func (c *Client) doRequest(httpMethod string, uri string, request Request, response interface{}) ([]byte, error) {
	return c.doRequestIter(httpMethod, uri, request, response, 0)
}

func toErrorWithNewlines(uri string, rawBody []byte) string {
	return fmt.Sprintf("Request URI was: %s Error message body: %s", uri, strings.ReplaceAll(string(rawBody), `\n`, "\n"))
}
