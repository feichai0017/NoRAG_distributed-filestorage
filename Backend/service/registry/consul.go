package registry

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
	"sync"
)

var (
	consulRegistry registry.Registry
	once           sync.Once
)

// GetConsulRegistry returns a singleton instance of Consul registry
func GetConsulRegistry() registry.Registry {
	once.Do(func() {
		consulRegistry = consul.NewRegistry(registry.Addrs("localhost:8500"))
	})
	return consulRegistry
}
