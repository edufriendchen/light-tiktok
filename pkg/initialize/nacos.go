package initialize

import (
	"log"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
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
		*constant.NewServerConfig(consts.NACOS_ADDR, 8848),
	}
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "log",
		CacheDir:            "cache",
		LogLevel:            "debug",
		Username:            consts.NACOS_USERNAME,
		Password:            consts.NACOS_PASSWORD,
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

	global.NacosClient = cli

	// Create naming client for service discovery
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
	})
	if err != nil {
		klog.Fatal(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "TIKTOK_JSON",
		Group:  "DEFAULT_GROUP",
	})

	global.Config, err = viper.NewConfig("", content)
	if err != nil {
		log.Println("nacos pull config err:", err)
	}
	return cli, nil
}
