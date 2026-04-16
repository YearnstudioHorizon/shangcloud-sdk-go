package shangcloud

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserBasicInfo struct {
	UserId   int    `json:"uid"`      // 用户UID
	Nickname string `json:"nickname"` // 用户昵称
	Mail     string `json:"mail"`     // 电子邮箱
	Avatar   string `json:"avatar"`   // 头像的URL
}
