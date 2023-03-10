package main

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user/userservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	registry_nacos "github.com/kitex-contrib/registry-nacos/registry"
	"gorm.io/plugin/opentelemetry/provider"
)

func Init() {
	initialize.InitNacos()
	initialize.InitCacheManager()
	initialize.InitNeo4j()
	initialize.InitJWT()
	initialize.InitSonyflake()
	klog.SetLogger(logrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
}

func main() {
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.USER_SERVICE_ADDR)
	if err != nil {
		panic(err)
	}
	Init()
	ctx := context.Background()
	defer global.Neo4jDriver.Close(ctx)

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.USER_SERVICE_NAME),
		provider.WithExportEndpoint("localhost:5775"),
		provider.WithInsecure(),
	)

	defer p.Shutdown(context.Background())
	svr := userservice.NewServer(new(UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry_nacos.NewNacosRegistry(global.NacosClient)),
		server.WithRegistryInfo(&registry.Info{
			ServiceName: consts.USER_SERVICE_NAME,
			Addr:        addr,
			Weight:      10,
			Tags:        nil,
		}),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.USER_SERVICE_NAME}),
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
