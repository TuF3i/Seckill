# Nacos 接入教程

## 初始化 Client

```go
import n "seckill/infrastructures/nacos"

client, err := n.NewNacosClient(
    n.WithHost("127.0.0.1"),
    n.WithPort(8848),
    n.WithUserName("admin"),
    n.WithPassword("admin"),
    n.WithNamespaceID("public"),
)
```

`Client` 结构体包含两个核心组件：

| 字段 | 类型 | 用途 |
|------|------|------|
| `NamingClient` | `naming_client.INamingClient` | 服务注册/发现 |
| `ConfigClient` | `config_client.IConfigClient` | 配置管理 |

---

## Kitex 注册中心接入

### 服务端注册

```go
import (
    nacosRegistry "github.com/kitex-contrib/registry-nacos/registry"
    n "seckill/infrastructures/nacos"
)

client, _ := n.NewNacosClient()

svr := xxxserver.NewServer(
    new(XxxServiceImpl),
    server.WithRegistry(nacosRegistry.NewNacosRegistry(client.NamingClient)),
)
svr.Run()
```

### 客户端发现

```go
import (
    nacosResolver "github.com/kitex-contrib/registry-nacos/resolver"
    n "seckill/infrastructures/nacos"
)

client, _ := n.NewNacosClient()

cli, err := xxxservice.MustNewClient(
    "target-service",
    client.WithResolver(nacosResolver.NewNacosResolver(client.NamingClient)),
)
```

---

## 配置中心接入

### 读取配置

```go
import (
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
    n "seckill/infrastructures/nacos"
)

client, _ := n.NewNacosClient()

content, err := client.ConfigClient.GetConfig(vo.ConfigParam{
    DataId: "seckill",
    Group:  "REDROCK",
})
```

### 配置加载器

项目在 `pkg/config/loader.go` 中封装了配置加载器，自动从 Nacos 拉取 JSON 配置并解析到 `configs.Config` 结构体：

```go
import (
    n "seckill/infrastructures/nacos"
    "seckill/pkg/config"
)

nacosClient, _ := n.NewNacosClient()
loader, _ := config.NewLoader(nacosClient, "seckill", "REDROCK")
cfg := loader.GetConfig()
// cfg.PostgreSQL.Host, cfg.Redis.MasterName, etc.
```
