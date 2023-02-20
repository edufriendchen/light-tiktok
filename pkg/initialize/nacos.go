package initialize

import (
	"log"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/edufriendchen/light-tiktok/pkg/viper"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// InitHertzNacos to init nacos
func InitNacos() (naming_client.INamingClient, error) {
	// Read configuration information from nacos
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
		Username:            "nacos",
		Password:            "nacosfriend0429",
	}
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		klog.Infof("Nacos Init Error: %v", err)
	}

	// Create naming client for service discovery
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
	})
	if err != nil {
		klog.Fatal(err)
	}

	//读取文件
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "TIKTOK_JSON",   //此处对应之前的网页配置的名称
		Group:  "DEFAULT_GROUP", //此处对应之前的网页配置的分组
	})

	global.Config, err = viper.NewConfig("", content)
	if err != nil {
		log.Println("config:", global.Config, "err:", err)
	}
	return cli, nil
}
