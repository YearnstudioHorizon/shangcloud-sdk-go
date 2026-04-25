# ShangCloud SDK for Go

 Go 语言 SDK，封装了授权登录与基础用户信息接口。

- Module: `github.com/YearnstudioHorizon/shangcloud-sdk-go`
- Package: `shangcloud`
- Go 版本: 1.25.0+
- License: [MIT](./LICENSE)

## 安装

```bash
go get github.com/YearnstudioHorizon/shangcloud-sdk-go
```

## 快速开始

下面是一个完整的 OAuth 授权码模式 (Authorization Code) 流程示例。

```go
package main

import (
    "fmt"
    "net/http"

    shangcloud "github.com/YearnstudioHorizon/shangcloud-sdk-go"
)

func main() {
    client := shangcloud.InitClient(
        "your-client-id",
        "your-client-secret",
        "https://your-app.example.com/oauth/callback",
    )

    // 生成授权跳转 URL，将用户引导到授权页
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, client.GenerateOAuthUrl(), http.StatusFound)
    })

    // 处理授权回调，使用 code 换取 User 实例
    http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
        code := r.URL.Query().Get("code")
        state := r.URL.Query().Get("state")

        user, err := client.GenerateUserInstance(code, state)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        // 拉取用户基本信息
        info, err := user.GetBasicInfo()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }

        fmt.Fprintf(w, "Hello, %s (uid=%d)\n", info.Nickname, info.UserId)
    })

    http.ListenAndServe(":8080", nil)
}
```

## 核心 API

### `InitClient(clientId, clientSecret, redirectUri string) *Client`

创建一个 SDK 客户端。默认 `Scope` 为 `user:basic`，`BaseUrl` 为 `https://api.yearnstudio.cn`，并使用内置的内存 KV 作为 state 存储。如需自定义这些字段，可以在返回的 `*Client` 上直接覆盖。

### `(*Client) GenerateOAuthUrl() string`

生成授权跳转 URL，内部会随机生成 state 并写入 `KvStorage`，用于后续回调校验。

### `(*Client) GenerateUserInstance(code, state string) (User, error)`

校验 state，使用授权码向 `/oauth/token` 换取 access token / refresh token，并返回实现了 `User` 接口的实例。

### `User` 接口

```go
type User interface {
    InitUser(accessToken, refreshToken, tokenType string, expiresIn int, c *Client)
    Save()
    IsExpired() bool
    GetBasicInfo() (UserBasicInfo, error)
    GetVariable(key string) (string, error)
    SetVariable(key, value string) error
    DeleteVariable(key string) error
}
```

SDK 提供了默认的内存实现 `UserInstance`。`IsExpired` 会提前 60 秒判断过期，`GetBasicInfo` 会请求 `/api/user/info`。

### 用户变量读写

`GetVariable` / `SetVariable` / `DeleteVariable` 通过 `/api/varibles` 操作当前用户的变量存储，需要授权时携带 `var:io` scope。

```go
if err := user.SetVariable("theme", "dark"); err != nil {
    // 处理错误
}
value, err := user.GetVariable("theme")
_ = user.DeleteVariable("theme")
```

### `UserBasicInfo`

```go
type UserBasicInfo struct {
    UserId   int    `json:"uid"`
    Nickname string `json:"nickname"`
    Mail     string `json:"mail"`
    Avatar   string `json:"avatar"`
}
```

## 自定义扩展

### 自定义 state 存储

`Client.KvStorage` 字段实现以下接口，将其替换为 Redis / 数据库等共享存储即可在多实例部署中复用：

```go
type TempVarStorage interface {
    SetTempVarible(key, value string)
    GetTempVarible(key string) (string, error)
    DeleteTempVarible(key string)
}
```

要求实现自行保证线程安全。

### 自定义 User 持久化

默认的 `UserInstance` 不会持久化 token。如果需要将 token 落库 / 写入会话，可以实现自己的 `User` 类型，在 `InitUser` / `Save` 中加入持久化逻辑，并自行调用底层接口完成令牌交换。

### 更换 ClientSecret

```go
client.SetClientSecret("new-secret")
```

## 注意事项

- `Client` 与 `UserInstance` 内嵌了 `noCopy`，请始终通过指针传递，避免值拷贝导致 `sync.Map` 等内部状态出现问题。
- `clientSecret` 与 token 字段为非导出字段，序列化时不会泄漏，但请仍然避免将 `Client` / `UserInstance` 直接打印到日志中。

## License

[MIT](./LICENSE)
