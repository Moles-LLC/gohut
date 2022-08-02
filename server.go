package gohut

import (
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

type Server struct {
	Banner struct {
		Image string `json:"image"`
		Tint  string `json:"tint"`
	} `json:"banner"`
	Categories       []string  `json:"categories"`
	ConnectedServers []string  `json:"connectedServers"`
	ID               string    `json:"id"`
	Icon             string    `json:"icon"`
	MOTD             string    `json:"motd"`
	MaxPlayers       int64     `json:"maxPlayers"`
	Name             string    `json:"name"`
	Online           bool      `json:"online"`
	Plan             string    `json:"plan"`
	PlayerCount      int64     `json:"playerCount"`
	Started          time.Time `json:"started"`
	Software         struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"versionType"`
	Suspended bool `json:"suspended"`
	Visible   bool `json:"visible"`
}

func parseServer(g gjson.Result) *Server {
	// I know this is stupid, but I cba making more structs for json parsing
	return &Server{
		Banner: struct {
			Image string `json:"image"`
			Tint  string `json:"tint"`
		}{
			Image: g.Get("default_banner_image").String(),
			Tint:  g.Get("default_banner_tint").String(),
		},
		Categories:       gjsonStringArray(g.Get("categories")),
		ConnectedServers: gjsonStringArray(g.Get("connectedServers")),
		ID:               g.Get("_id").String(),
		Icon:             g.Get("icon").String(),
		MOTD:             g.Get("motd").String(),
		MaxPlayers:       g.Get("maxPlayers").Int(),
		Name:             g.Get("name").String(),
		Online:           g.Get("online").Bool(),
		Plan:             g.Get("activeServerPlan").String(),
		PlayerCount:      g.Get("playerCount").Int(),
		Software: struct {
			Type    string `json:"type"`
			Version string `json:"version"`
		}{
			Type: g.Get("server_version_type").String(),
		},
		Suspended: g.Get("suspended").Bool(),
		Visible:   g.Get("visibility").Bool(),
	}
}

func (c *Client) GetServerByName(name string) (*Server, error) {
	return c.getServer(name, true)
}

func (c *Client) GetServerByID(id string) (*Server, error) {
	return c.getServer(id, false)
}

func (c *Client) getServer(id string, name bool) (*Server, error) {
	url := c.BaseUrl + "/server/" + id
	if name {
		url += "?byName=true"
	}

	body, err := c.MakeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	g := gjson.ParseBytes(body)

	return parseServer(g.Get("server")), nil
}
