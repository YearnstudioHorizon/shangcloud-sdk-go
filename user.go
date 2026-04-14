package shangcloud

type User interface {
	InitUser(string, string, int, *Client) // 在实例化后会被立即调用, 传入accessToken refreshToken userId及Client指针
	Save()                                 // 在数据变更后会被调用
}

type UserInstance struct {
	accessToken  string  `json:"-"`
	refreshToken string  `json:"-"`
	UserId       int     `json:"user_id"`
	client       *Client `json:"-"`
}

// 由于是内存存储, 不需要具体实现
func (user *UserInstance) Save() {}

func (user *UserInstance) InitUser(accessToken string, refreshToken string, userId int, client *Client) {
	user.accessToken = accessToken
	user.refreshToken = refreshToken
	user.UserId = userId
	user.client = client
}

var _ User = &UserInstance{} // 测试接口是否实现
