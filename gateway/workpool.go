package gateway

import (
	"fmt"
	"github.com/panjf2000/ants"
	"plato/common/config"
)

var wPool *ants.Pool

func initWorkPool() {
	var err error
	if wPool, err = ants.NewPool(config.GetGatewayWorkerPoolNum()); err != nil {
		fmt.Printf("InitWorkPoll.err :%s num:%d\n", err.Error(), config.GetGatewayWorkerPoolNum())
	}
}
