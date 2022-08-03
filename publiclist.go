package gohut

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func parseServerInList(g gjson.Result) *Server {
	// I know this is stupid, but I cba making more structs for json parsing
	return &Server{
		Banner: struct {
			Image string `json:"image"`
			Tint  string `json:"tint"`
		}{
			Image: g.Get("default_banner_image").String(),
			Tint:  g.Get("default_banner_tint").String(),
		},
		Categories:       gjsonStringArray(g.Get("allCategories")),
		ConnectedServers: gjsonStringArray(g.Get("staticInfo.connectedServers")),
		ID:               g.Get("staticInfo._id").String(),
		Icon:             g.Get("icon").String(),
		MOTD:             g.Get("motd").String(),
		MaxPlayers:       g.Get("maxPlayers").Int(),
		Name:             g.Get("name").String(),
		Plan:             g.Get("staticInfo.serverPlan").String(),
		PlayerCount:      g.Get("playerData.playerCount").Int(),
		Started:          time.Unix(g.Get("staticInfo.serviceStartDate").Int(), 0),
		Software: struct {
			Type    string `json:"type"`
			Version string `json:"version"`
		}{
			Type:    g.Get("server_version.type").String(),
			Version: g.Get("server_version.version").String(),
		},
		Suspended: false,                      // if it's in the server list, it's most likely not suspended
		Visible:   g.Get("visibility").Bool(), // if it's in the server list, it's most likely visible, but we still wanna check just to be 100% sure
	}
}

func (c *Client) GetPublicServerList(query string, limit int, offset int) ([]*Server, int, error) {
	body, err := c.MakeRequest(http.MethodGet, c.BaseUrl+"/servers?q="+query+"&limit="+fmt.Sprintf("%d", limit)+"&offset="+fmt.Sprintf("%d", offset), nil)
	if err != nil {
		return nil, 0, err
	}

	g := gjson.ParseBytes(body)
	serverArr := g.Get("servers").Array()
	servers := make([]*Server, len(serverArr))

	for i, rawServer := range g.Get("servers").Array() {
		server := parseServerInList(rawServer)
		servers[i] = server
	}

	total := int(g.Get("total_search_results").Int())

	return servers, total, nil
}

func (c *Client) GetPublicServerListDetails(query string) ([]*Server, error) {
	firstServers, totalServers, err := c.GetPublicServerList(query, 100, 0)
	if err != nil {
		return nil, err
	}

	servers := make([]*Server, totalServers)
	serverIndex := 0

	for _, server := range firstServers {
		servers[serverIndex] = server
		serverIndex++
	}

	if totalServers > 100 {
		channel := make(chan *Server)

		for i := 100; i < totalServers; i += 100 {
			go func(i int) {
				goServers, _, err := c.GetPublicServerList(query, 100, i)

				defer func() {
					channel <- &Server{
						ID: "\x00",
					}
				}()

				if err != nil {
					// too bad
					return
				}

				for _, server := range goServers {
					channel <- server
				}
			}(i)
		}

		tasks := totalServers / 100

		for tasks > 0 {
			server := <-channel

			// a thread has finished
			if server.ID == "\x00" {
				tasks--
				continue
			}

			if serverIndex >= totalServers {
				// if a server goes up while fetching
				servers = append(servers, server)
			} else {
				servers[serverIndex] = server
			}

			serverIndex++
		}
	}

	// if a server goes down while fetching
	servers = servers[:serverIndex]

	return servers, nil
}
