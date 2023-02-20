// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/edufriendchen/httpvlog"
	"github.com/edufriendchen/light-tiktok/cmd/api/biz/rpc"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/registry/nacos"
)

func Init() {
	global.NacosClient, _ = initialize.InitNacos()
	rpc.Init()
	initialize.InitJWT()
	hlog.SetLogger(logrus.NewLogger())
	hlog.SetLevel(hlog.LevelInfo)
}

func main() {
	Init()
	tracer, cfg := tracing.NewServerTracer()
	h := server.New(
		server.WithHostPorts(consts.API_SERVICE_ADDR),
		server.WithHandleMethodNotAllowed(true),
		server.WithRegistry(nacos.NewNacosRegistry(global.NacosClient), &registry.Info{
			ServiceName: consts.API_SERVICE_NAME,
			Addr:        utils.NewNetAddr(consts.TCP, consts.API_SERVICE_ADDR),
			Weight:      10,
			Tags:        nil,
		}),
		server.WithMaxRequestBodySize(100<<20),
		tracer,
	)
	h.Use(tracing.ServerMiddleware(cfg))
	h.Use(httpvlog.Logger())
	register(h)
	h.Spin()
}
