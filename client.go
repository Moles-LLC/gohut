package gohut

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	UserAgent  string
	BaseUrl    string
}

func NewClient() *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		UserAgent:  "Samsung Smart Fridge",
		BaseUrl:    "https://api.minehut.com",
	}
}

func (c *Client) newRequest(method, url string, body interface{}) (*http.Request, error) {
	var rawBody []byte
	var err error

	if body != nil {
		rawBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) MakeRequest(method, url string, body interface{}) ([]byte, error) {
	req, err := c.newRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if 200 < res.StatusCode || res.StatusCode > 299 {
		return nil, errors.New("status code not 2xx - " + fmt.Sprint(res.StatusCode))
	}

	return io.ReadAll(res.Body)
}
