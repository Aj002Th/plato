package ipconfig

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"plato/common/config"
	"plato/ipconfig/domain"
	"plato/ipconfig/source"
)

func RunMain(path string) {
	config.Init(path) // 初始化配置
	source.Init()     // 初始化数据源连接层(etcd)
	domain.Init()     // 初始化 ip 调度层

	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpList)
	s.Spin()
}
