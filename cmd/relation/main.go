package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/edufriendchen/light-tiktok/kitex_gen/relation/relationservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	registry_nacos "github.com/kitex-contrib/registry-nacos/registry"
)

// User RPC Server 端配置初始化
func Init() {
	initialize.InitJWT()
	initialize.InitNeo4j()
	klog.SetLogger(logrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
}

func main() {
	Init()
	cli, _ := initialize.InitNacos()
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.RELATION_SERVICE_ADDR)
	if err != nil {
		panic(err)
	}
	svr := relationservice.NewServer(new(RelationServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry_nacos.NewNacosRegistry(cli)),
		server.WithRegistryInfo(&registry.Info{
			ServiceName: consts.RELATION_SERVICE_NAME,
			Addr:        addr,
			Weight:      10,
			Tags:        nil,
		}),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.RELATION_SERVICE_NAME}),
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
