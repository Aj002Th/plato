package discovery

import (
	"context"
	"plato/common/config"
	"testing"
	"time"
)

func TestServiceDiscovery(t *testing.T) {
	config.Init("../../plato.yaml")
	ctx := context.Background()
	ser := NewServiceDiscovery(&ctx)
	defer ser.Close()
	ser.Watch("/web/", func(key, value string) {}, func(key, value string) {})
	ser.Watch("/gRPC/", func(key, value string) {}, func(key, value string) {})
	for {
		select {
		case <-time.Tick(10 * time.Second):
		}
	}
}
