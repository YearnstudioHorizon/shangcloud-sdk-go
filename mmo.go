package shangcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MMONewRoomResponse struct {
	ConnectKey string `json:"connect_key"`
	EdgeUrl    string `json:"edge_url"`
	RoomId     string `json:"room_id"`
	Protocol   string `json:"protocol"`
}

type MMOJoinRoomResponse struct {
	ConnectKey  string `json:"connect_key"`
	EdgeUrl     string `json:"edge_url"`
	RoomId      string `json:"room_id"`
	Protocol    string `json:"protocol"`
	AssignedUid string `json:"assigned_uid,omitempty"`
}

func (c *Client) mmoRequest(path string, body any, accessToken string, tokenType string, roomId string, protocol string) ([]byte, error) {
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, path)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	if roomId != "" {
		req.Header.Set("X-MMO-Room", roomId)
	}
	if protocol != "" {
		req.Header.Set("X-MMO-Protoctl", protocol)
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return data, fmt.Errorf("server returned error status: %d, body: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// NewRoom creates a new MMO room. The caller automatically becomes the room owner.
// protocol: "tcp" or "websocket"; empty string uses server default.
func (user *UserInstance) NewRoom(protocol string) (MMONewRoomResponse, error) {
	data, err := user.client.mmoRequest("/api/mmo/room/new", struct{}{}, user.accessToken, user.TokenType, "", protocol)
	if err != nil {
		return MMONewRoomResponse{}, err
	}
	var resp MMONewRoomResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return MMONewRoomResponse{}, err
	}
	return resp, nil
}

// JoinRoom joins an existing MMO room.
// protocol: "tcp" or "websocket"; empty string uses server default.
func (user *UserInstance) JoinRoom(roomId string, protocol string) (MMOJoinRoomResponse, error) {
	data, err := user.client.mmoRequest("/api/mmo/room/join", struct{}{}, user.accessToken, user.TokenType, roomId, protocol)
	if err != nil {
		return MMOJoinRoomResponse{}, err
	}
	var resp MMOJoinRoomResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return MMOJoinRoomResponse{}, err
	}
	return resp, nil
}

// SetRoomConfig sets room configuration. Only the room owner can call this.
func (user *UserInstance) SetRoomConfig(roomId string, allowMultiLogin bool) error {
	_, err := user.client.mmoRequest("/api/mmo/room/config", map[string]any{"allow_multi_login": allowMultiLogin}, user.accessToken, user.TokenType, roomId, "")
	return err
}

// SetRoomData sets a key-value pair in the room's extra data. Only the room owner can call this.
// dataType: "number", "string", or "boolean"; empty string uses server default ("string").
func (user *UserInstance) SetRoomData(roomId string, key string, value any, dataType string) error {
	body := map[string]any{"key": key, "value": value}
	if dataType != "" {
		body["type"] = dataType
	}
	_, err := user.client.mmoRequest("/api/mmo/room/data/set", body, user.accessToken, user.TokenType, roomId, "")
	return err
}

// GetRoomData retrieves all extra data stored in the room.
func (user *UserInstance) GetRoomData(roomId string) (map[string]any, error) {
	data, err := user.client.mmoRequest("/api/mmo/room/data/get", struct{}{}, user.accessToken, user.TokenType, roomId, "")
	if err != nil {
		return nil, err
	}
	var resp struct {
		ExtraData map[string]any `json:"extra_data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.ExtraData, nil
}

// DeleteRoomData deletes a key from the room's extra data. Only the room owner can call this.
func (user *UserInstance) DeleteRoomData(roomId string, key string) error {
	_, err := user.client.mmoRequest("/api/mmo/room/data/delete", map[string]any{"key": key}, user.accessToken, user.TokenType, roomId, "")
	return err
}

// KickUser forcibly removes a user from the room. Only the room owner can call this.
// The room owner cannot kick themselves.
func (user *UserInstance) KickUser(roomId string, targetUid string) error {
	_, err := user.client.mmoRequest("/api/mmo/room/kick", map[string]any{"target_uid": targetUid}, user.accessToken, user.TokenType, roomId, "")
	return err
}

// GetRoomUserCount returns the current number of users in the room.
func (user *UserInstance) GetRoomUserCount(roomId string) (int, error) {
	data, err := user.client.mmoRequest("/api/mmo/room/usercount", struct{}{}, user.accessToken, user.TokenType, roomId, "")
	if err != nil {
		return 0, err
	}
	var resp struct {
		UserCount int `json:"user_count"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, err
	}
	return resp.UserCount, nil
}
