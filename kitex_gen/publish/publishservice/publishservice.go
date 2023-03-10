// Code generated by Kitex v0.4.4. DO NOT EDIT.

package publishservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	publish "github.com/edufriendchen/light-tiktok/kitex_gen/publish"
)

func serviceInfo() *kitex.ServiceInfo {
	return publishServiceServiceInfo
}

var publishServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "PublishService"
	handlerType := (*publish.PublishService)(nil)
	methods := map[string]kitex.MethodInfo{
		"ActionPulish":    kitex.NewMethodInfo(actionPulishHandler, newPublishServiceActionPulishArgs, newPublishServiceActionPulishResult, false),
		"MGetPublishList": kitex.NewMethodInfo(mGetPublishListHandler, newPublishServiceMGetPublishListArgs, newPublishServiceMGetPublishListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "publish",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func actionPulishHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*publish.PublishServiceActionPulishArgs)
	realResult := result.(*publish.PublishServiceActionPulishResult)
	success, err := handler.(publish.PublishService).ActionPulish(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newPublishServiceActionPulishArgs() interface{} {
	return publish.NewPublishServiceActionPulishArgs()
}

func newPublishServiceActionPulishResult() interface{} {
	return publish.NewPublishServiceActionPulishResult()
}

func mGetPublishListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*publish.PublishServiceMGetPublishListArgs)
	realResult := result.(*publish.PublishServiceMGetPublishListResult)
	success, err := handler.(publish.PublishService).MGetPublishList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newPublishServiceMGetPublishListArgs() interface{} {
	return publish.NewPublishServiceMGetPublishListArgs()
}

func newPublishServiceMGetPublishListResult() interface{} {
	return publish.NewPublishServiceMGetPublishListResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) ActionPulish(ctx context.Context, req *publish.ActionRequest) (r *publish.ActionResponse, err error) {
	var _args publish.PublishServiceActionPulishArgs
	_args.Req = req
	var _result publish.PublishServiceActionPulishResult
	if err = p.c.Call(ctx, "ActionPulish", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) MGetPublishList(ctx context.Context, req *publish.PublishRequest) (r *publish.PublishResponse, err error) {
	var _args publish.PublishServiceMGetPublishListArgs
	_args.Req = req
	var _result publish.PublishServiceMGetPublishListResult
	if err = p.c.Call(ctx, "MGetPublishList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
