package shangcloud

import (
	"fmt"
	"net/url"
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

	url := fmt.Sprintf("%v/oauth/authorize%v", c.BaseUrl, params.Encode())
	return url
}

func (c *Client) SetClientSecret(clientSecret string) {
	c.clientSecret = clientSecret
}
