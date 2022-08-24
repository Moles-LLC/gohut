package gohut

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var DefaultClient = NewClient()

type Client struct {
	HttpClient *http.Client
	UserAgent  string
	BaseUrl    string

	AccessToken string
	ProfileID   string
	SessionID   string
}

func NewClient() *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		UserAgent:  "Samsung Smart Fridge",
		BaseUrl:    "https://api.minehut.com",
	}
}

func (c *Client) newRequest(method, url string, body interface{}, authenticate bool) (*http.Request, error) {
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

	if authenticate {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
		req.Header.Add("x-profile-id", c.ProfileID)
		req.Header.Add("x-session-id", c.SessionID)
	}

	return req, nil
}

func (c *Client) MakeRequest(method, url string, body any, authenticate bool) ([]byte, error) {
	req, err := c.newRequest(method, url, body, authenticate)
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
