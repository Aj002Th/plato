package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	Init("../../plato.yaml")
	fmt.Printf("GetEndpointsForDiscovery: %v\n", GetEndpointsForDiscovery())
	fmt.Printf("GetTimeoutForDiscovery: %v\n", GetTimeoutForDiscovery())
	fmt.Printf("GetServicePathForIPConf: %v\n", GetServicePathForIPConf())
	fmt.Printf("IsDebug: %v\n", IsDebug())
}
