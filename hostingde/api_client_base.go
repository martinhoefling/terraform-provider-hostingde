package hostingde

// Adapted from hostingde provider in https://github.com/go-acme/lego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	authToken      string
	ownerAccountId string
	HTTPClient     *http.Client
}

func NewClient(authToken string, ownerAccountId string) *Client {
	c := Client{authToken: authToken, ownerAccountId: ownerAccountId, HTTPClient: &http.Client{}}
	return &c
}

func (c *Client) post(uri string, request Request, response interface{}) ([]byte, error) {
	if request.getAuthToken() == "" {
		request.setAuthToken(c.authToken)
	}
	if request.getOwnerAccountId() == "" {
		request.setOwnerAccountId(c.ownerAccountId)
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	log.Print(string(body))
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying API: %v", err)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(toUnreadableBodyMessage(uri, content))
	}

	err = json.Unmarshal(content, response)
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, toUnreadableBodyMessage(uri, content))
	}

	return content, nil
}

func toUnreadableBodyMessage(uri string, rawBody []byte) string {
	return fmt.Sprintf("the request %s sent a response with a body which is an invalid format:", uri, strings.Replace(string(rawBody), `\n`, "\n", -1))
}
