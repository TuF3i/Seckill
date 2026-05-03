# Nacos 接入教程

## 初始化 Client

```go
import n "seckill/infrastructures/nacos"

client, err := n.NewNacosClient(
    n.WithHost("127.0.0.1"),
    n.WithPort(8848),
    n.WithUserName("nacos"),
    n.WithPassword("nacos"),
    n.WithNamespaceID("dev"),
)
```

`Client` 结构体包含两个核心组件：

| 字段 | 类型 | 用途 |
|------|------|------|
| `NamingClient` | `naming_client.INamingClient` | 服务注册/发现 |
| `ConfigClient` | `config_client.IConfigClient` | 配置管理 |

---

## Kitex 注册中心接入

### 安装依赖

```bash
go get github.com/kitex-contrib/registry-nacos
```

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

## Viper 配置中心接入

### 安装依赖

```bash
go get github.com/nacos-group/nacos-sdk-go/v2/clients/config_client
```

### 远程配置监听

```go
import (
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
    n "seckill/infrastructures/nacos"
)

client, _ := n.NewNacosClient()

// 读取配置
content, err := client.ConfigClient.GetConfig(vo.ConfigParam{
    DataId: "app-config",
    Group:  "DEFAULT_GROUP",
})

// 监听配置变更
client.ConfigClient.ListenConfig(vo.ConfigParam{
    DataId: "app-config",
    Group:  "DEFAULT_GROUP",
    OnChange: func(namespace, group, dataId, data string) {
        // 配置变更回调
    },
})
```

### 与 Viper 结合

```go
import (
    "github.com/spf13/viper"
    n "seckill/infrastructures/nacos"
)

client, _ := n.NewNacosClient()

content, _ := client.ConfigClient.GetConfig(vo.ConfigParam{
    DataId: "app.yaml",
    Group:  "DEFAULT_GROUP",
})

viper.SetConfigType("yaml")
viper.ReadConfig(strings.NewReader(content))
```
