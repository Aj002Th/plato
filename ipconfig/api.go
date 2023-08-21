package ipconfig

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"plato/ipconfig/domain"
)

type Response struct {
	Code    int
	Message string
	Data    any
}

// GetIpList 获得 gateway 服务 ip
func GetIpList(c context.Context, ctx *app.RequestContext) {
	// 兜底
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()

	// 1.构造请求信息
	ipCtx := domain.NewIpConfigContext(&c, ctx)
	// 2.进行ip调度计算
	eds := domain.Dispatch(ipCtx)
	// 3.取得分 top5 作为结果返回
	if len(eds) > 5 {
		eds = eds[:5]
	}
	ipCtx.AppCtx.JSON(consts.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    eds,
	})
}
