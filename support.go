package gohut

import "net/http"

type SupportRequestBody struct {
	UserID            string `json:"user_id"`
	UserEmail         string `json:"user_email"`
	UserMinecraftName string `json:"user_minecraft_name"`
	ProblemType       string `json:"problem_type"`
	ServerID          string `json:"server_id"`
	ServerName        string `json:"server_name"`
	ServerVersion     string `json:"server_version"`
	ProblemDetails    string `json:"problem_details"`
	MinecraftType     string `json:"minecraft_type"`
	MinecraftVersion  string `json:"minecraft_version"`
	MinecraftMods     string `json:"minecraft_mods"`
}

func (c *Client) SendSupportRequest(body SupportRequestBody) error {
	_, err := c.MakeRequest(http.MethodPost, c.BaseUrl+"/new/support", body, true)
	return err
}
