package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/edufriendchen/light-tiktok/kitex_gen/message/messageservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	registry_nacos "github.com/kitex-contrib/registry-nacos/registry"
)

func Init() {
	initialize.InitDB()
	initialize.InitJWT()
	klog.SetLogger(logrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
}

func main() {
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.MessageServiceAddr)
	if err != nil {
		panic(err)
	}
	Init()
	cli, _ := initialize.InitNacos()
	// provider.NewOpenTelemetryProvider(
	// 	provider.WithServiceName(consts.UserServiceName),
	// 	provider.WithExportEndpoint(consts.ExportEndpoint),
	// 	provider.WithInsecure(),
	// )
	svr := messageservice.NewServer(new(MessageServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry_nacos.NewNacosRegistry(cli)),
		server.WithRegistryInfo(&registry.Info{
			ServiceName: consts.MessageServiceName,
			Addr:        addr,
			Weight:      10,
			Tags:        nil,
		}),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.MessageServiceName}),
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
