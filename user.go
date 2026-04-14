package shangcloud

import (
	"time"
)

type User interface {
	InitUser(string, string, int, *Client) // 在实例化后会被立即调用
	Save()                                 // 在数据变更后会被调用
	IsExpired() bool                       // 检查 Token 是否过期
}

type UserInstance struct {
	accessToken  string    `json:"-"`
	refreshToken string    `json:"-"`
	ExpiresIn    int       `json:"expires_in"`  // 原始有效时长（秒）
	ExpiryTime   time.Time `json:"expiry_time"` // 具体的过期时刻
	client       *Client   `json:"-"`
	noCopy                 // 禁止拷贝
}

// 由于是内存存储, 不需要具体实现
func (user *UserInstance) Save() {}

func (user *UserInstance) InitUser(accessToken string, refreshToken string, expiresIn int, client *Client) {
	user.accessToken = accessToken
	user.refreshToken = refreshToken
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

var _ User = &UserInstance{} // 测试接口是否实现
