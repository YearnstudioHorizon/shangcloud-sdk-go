package shangcloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	ClientId     string
	clientSecret string
	RedirectUri  string
	Scope        string
	BaseUrl      string
	KvStorage    TempVarStorage
	noCopy       // 禁止拷贝
}

func InitClient(clientId string, clientSecret string, redirectUri string) *Client {
	storage := newRamKv()
	return &Client{
		ClientId:     clientId,
		clientSecret: clientSecret,
		RedirectUri:  redirectUri,
		Scope:        "user:basic",
		BaseUrl:      "https://api.yearnstudio.cn",
		KvStorage:    storage,
	}
}

// 生成OAuth授权跳转的URL
func (c *Client) GenerateOAuthUrl() string {
	// 生成state
	state := generateRandomString(10)

	// 存入storage
	c.KvStorage.SetTempVarible(state, "0")

	// 构造param参数
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("state", state)
	params.Add("client_id", c.ClientId)
	params.Add("redirect_uri", c.RedirectUri)
	params.Add("scope", c.Scope)

	url := fmt.Sprintf("%s/oauth/authorize%s", c.BaseUrl, params.Encode())
	return url
}

// 更换ClientSecret
func (c *Client) SetClientSecret(clientSecret string) {
	c.clientSecret = clientSecret
}

// 基于ClientId及ClientSecret生成Authorization头
func (c *Client) generateAuthorizeHeader() string {
	raw := fmt.Sprintf("%s:%s", c.ClientId, c.clientSecret)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

// 使用内置的UserInstace生成User接口实例
func (c *Client) GenerateUserInstance(code string, state string) (User, error) {
	// 查询state是否存在
	_, err := c.KvStorage.GetTempVarible(state)
	if err != nil {
		return nil, err
	}
	// 移除state
	c.KvStorage.DeleteTempVarible(state)

	// 构造请求参数
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectUri)

	// 构造请求对象
	reqUrl := fmt.Sprintf("%s/oauth/token", c.BaseUrl)
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// 设置 Header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+c.generateAuthorizeHeader()) // 使用工具函数生成

	// 发起请求
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth failed with status: %d", resp.StatusCode)
	}

	// 解析响应
	var tResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tResp); err != nil {
		return nil, err
	}

	// 生成user实例
	user := &UserInstance{}

	user.InitUser(tResp.AccessToken, tResp.RefreshToken, tResp.TokenType, tResp.ExpiresIn, c)
	return user, nil
}

func (c *Client) getUserBasicInfo(accessToken string, tokenType string) (UserBasicInfo, error) {
	body := map[string]string{}
	data, err := c.request("/api/user/info", body, accessToken, tokenType)
	if err != nil {
		return UserBasicInfo{}, err
	}
	var basicInfo UserBasicInfo
	json.Unmarshal(data, &basicInfo)
	return basicInfo, nil
}
