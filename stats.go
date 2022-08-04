package gohut

import (
	"net/http"

	"github.com/tidwall/gjson"
)

type SimpleStatsResponse struct {
	PlayerCount int64 `json:"playerCount"`
	ServerCount int64 `json:"serverCount"`
	ServerMax   int64 `json:"serverMax"`
	RamCount    int64 `json:"ramCount"`
	RamMax      int64 `json:"ramMax"`
}

func parseSimpleStatsResponse(g gjson.Result) *SimpleStatsResponse {
	return &SimpleStatsResponse{
		PlayerCount: g.Get("player_count").Int(),
		ServerCount: g.Get("server_count").Int(),
		ServerMax:   g.Get("server_max").Int(),
		RamCount:    g.Get("ram_count").Int(),
		RamMax:      g.Get("ram_max").Int(),
	}
}

func (c *Client) GetSimpleStats() (*SimpleStatsResponse, error) {
	body, err := c.MakeRequest(http.MethodGet, c.BaseUrl+"/network/simple_stats", nil)
	if err != nil {
		return nil, err
	}

	g := gjson.ParseBytes(body)

	return parseSimpleStatsResponse(g), nil
}

type PlayerDistributionResponse struct {
	Bedrock struct {
		Total        int64 `json:"total"`
		Lobby        int64 `json:"lobby"`
		PlayerServer int64 `json:"playerServer"`
	} `json:"bedrock"`
	Java struct {
		Total        int64 `json:"total"`
		Lobby        int64 `json:"lobby"`
		PlayerServer int64 `json:"playerServer"`
	} `json:"java"`
}

func parsePlayerDistributionResponse(g gjson.Result) *PlayerDistributionResponse {
	return &PlayerDistributionResponse{
		Bedrock: struct {
			Total        int64 `json:"total"`
			Lobby        int64 `json:"lobby"`
			PlayerServer int64 `json:"playerServer"`
		}{
			Total:        g.Get("bedrockTotal").Int(),
			Lobby:        g.Get("bedrockLobby").Int(),
			PlayerServer: g.Get("bedrockPlayerServer").Int(),
		},
		Java: struct {
			Total        int64 `json:"total"`
			Lobby        int64 `json:"lobby"`
			PlayerServer int64 `json:"playerServer"`
		}{
			Total:        g.Get("javaTotal").Int(),
			Lobby:        g.Get("javaLobby").Int(),
			PlayerServer: g.Get("javaPlayerServer").Int(),
		},
	}
}

func (c *Client) GetPlayerDistribution() (*PlayerDistributionResponse, error) {
	body, err := c.MakeRequest(http.MethodGet, c.BaseUrl+"/network/players/distribution", nil)
	if err != nil {
		return nil, err
	}

	g := gjson.ParseBytes(body)

	return parsePlayerDistributionResponse(g), nil
}
