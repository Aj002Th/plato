package config

import (
	"github.com/spf13/viper"
	"time"
)

// Init 读配置文件, 获取配置信息
func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// GetGlobalEnv 获取环境信息
func GetGlobalEnv() string {
	return viper.GetString("global.env")
}

// GetEndpointsForDiscovery 获取服务发现的地址
func GetEndpointsForDiscovery() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

// GetTimeoutForDiscovery 获取连接服务发现集群的超时时间 单位:秒
func GetTimeoutForDiscovery() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

// GetServicePathForIPConf 获取 ipconfig 在 etcd 中的路径
func GetServicePathForIPConf() string {
	return viper.GetString("ip_conf.service_path")
}

// IsDebug 判断是不是debug环境
func IsDebug() bool {
	env := viper.GetString("global.env")
	return env == "debug"
}
