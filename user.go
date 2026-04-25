package shangcloud

import (
	"time"
)

type User interface {
	InitUser(string, string, string, int, *Client) // 在实例化后会被立即调用
	Save()                                         // 在数据变更后会被调用
	IsExpired() bool                               // 检查 Token 是否过期
	GetBasicInfo() (UserBasicInfo, error)          // 生成UserBasicInfo结构体
	GetVariable(key string) (string, error)        // 读取用户变量
	SetVariable(key string, value string) error    // 写入用户变量
	DeleteVariable(key string) error               // 删除用户变量
}

type UserInstance struct {
	accessToken  string    `json:"-"`
	refreshToken string    `json:"-"`
	ExpiresIn    int       `json:"expires_in"`  // 原始有效时长（秒）
	ExpiryTime   time.Time `json:"expiry_time"` // 具体的过期时刻
	client       *Client   `json:"-"`
	TokenType    string    `json:"token_type"`
	noCopy                 // 禁止拷贝
}

// 由于是内存存储, 不需要具体实现
func (user *UserInstance) Save() {}

func (user *UserInstance) InitUser(accessToken string, refreshToken string, tokenType string, expiresIn int, client *Client) {
	user.accessToken = accessToken
	user.refreshToken = refreshToken
	user.TokenType = tokenType
	user.ExpiresIn = expiresIn
	user.client = client

	// 计算并存入过期时刻
	user.ExpiryTime = time.Now().Add(time.Duration(expiresIn) * time.Second)

	user.Save()
}

// IsExpired 检查当前 Token 是否已经失效
func (user *UserInstance) IsExpired() bool {
	// 提前60秒
	return time.Now().Add(60 * time.Second).After(user.ExpiryTime)
}

// 获取基础信息
func (user *UserInstance) GetBasicInfo() (UserBasicInfo, error) {
	return user.client.getUserBasicInfo(user.accessToken, user.TokenType)
}

// 读取用户变量
func (user *UserInstance) GetVariable(key string) (string, error) {
	return user.client.variableAction("read", key, "", user.accessToken, user.TokenType)
}

// 写入用户变量
func (user *UserInstance) SetVariable(key string, value string) error {
	_, err := user.client.variableAction("write", key, value, user.accessToken, user.TokenType)
	return err
}

// 删除用户变量
func (user *UserInstance) DeleteVariable(key string) error {
	_, err := user.client.variableAction("delete", key, "", user.accessToken, user.TokenType)
	return err
}

var _ User = &UserInstance{} // 测试接口是否实现
