package shangcloud

import (
	"encoding/json"
	"fmt"
)

// 变量读写接口的响应结构
type variableResponse struct {
	Value string `json:"value"`
	Error string `json:"error,omitempty"`
}

// 通过远端 /api/varibles 接口操作用户变量
func (c *Client) variableAction(action string, key string, value string, accessToken string, tokenType string) (string, error) {
	body := map[string]string{
		"key":    key,
		"action": action,
		"value":  value,
	}
	data, err := c.request("/api/varibles", body, accessToken, tokenType)
	if err != nil {
		return "", err
	}
	var resp variableResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("decode variable response failed: %w", err)
	}
	if resp.Error != "" {
		return "", fmt.Errorf("variable %s failed: %s", action, resp.Error)
	}
	return resp.Value, nil
}
