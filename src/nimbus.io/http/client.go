package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	credentials *Credentials
	httpClient  *http.Client
	baseAddress string
}

type Response struct {
	StatusCode int
	Status     string
	Body       []byte
}

func NewClient(credentials *Credentials, baseAddress string) *Client {
	return &Client{
		credentials,
		&http.Client{},
		baseAddress,
	}
}

func (client *Client) Request(method string, baseURI string) (*Response, error) {
	var err error

	current_time := time.Now()
	timestamp := current_time.Unix()
	uri := fmt.Sprintf("http://%s%s", client.baseAddress, baseURI)

	request, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	authString := ComputeAuthString(client.credentials, method, timestamp,
		baseURI)
	request.Header.Add("Authorization", authString)
	request.Header.Add("x-nimbus-io-timestamp", fmt.Sprintf("%d", timestamp))
	request.Header.Add("agent", "gonimbusio/1.0")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{response.StatusCode, response.Status, body}, nil
}