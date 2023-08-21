package domain

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type IpConfigContext struct {
	Ctx       *context.Context
	AppCtx    *app.RequestContext
	ClientCtx *ClientContext
}

// ClientContext 客户端请求时携带的额外信息
type ClientContext struct {
	IP string `json:"ip"`
}

func NewIpConfigContext(ctx *context.Context, appCtx *app.RequestContext) *IpConfigContext {
	return &IpConfigContext{
		Ctx:       ctx,
		AppCtx:    appCtx,
		ClientCtx: &ClientContext{},
	}
}
