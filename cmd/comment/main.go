package main

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	comment "github.com/edufriendchen/light-tiktok/kitex_gen/comment/commentservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	registry_nacos "github.com/kitex-contrib/registry-nacos/registry"
)

func Init() {
	initialize.InitDB()
	initialize.InitNeo4j()
	initialize.InitJWT()
	klog.SetLogger(logrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
}

func main() {
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.COMMENT_SERVICE_ADDR)
	if err != nil {
		panic(err)
	}
	Init()
	ctx := context.Background()
	defer global.Neo4jDriver.Close(ctx)
	cli, _ := initialize.InitNacos()
	// provider.NewOpenTelemetryProvider(
	// 	provider.WithServiceName(consts.UserServiceName),
	// 	provider.WithExportEndpoint(consts.ExportEndpoint),
	// 	provider.WithInsecure(),
	// )
	svr := comment.NewServer(new(CommentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry_nacos.NewNacosRegistry(cli)),
		server.WithRegistryInfo(&registry.Info{
			ServiceName: consts.COMMENT_SERVICE_NAME,
			Addr:        addr,
			Weight:      10,
			Tags:        nil,
		}),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.COMMENT_SERVICE_NAME}),
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
