package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"plato/common/config"
	"sync"
)

type ServiceDiscovery struct {
	cli  *clientv3.Client
	lock sync.Mutex
	ctx  *context.Context
}

// NewServiceDiscovery 创建服务发现
func NewServiceDiscovery(ctx *context.Context) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		panic(err)
	}

	return &ServiceDiscovery{
		cli:  cli,
		lock: sync.Mutex{},
		ctx:  ctx,
	}
}

// Watch 启动监听
func (s *ServiceDiscovery) Watch(prefix string, set, del func(k, v string)) error {
	// 1.依据 prefix 获得现有的所有 key, 记录信息
	resp, err := s.cli.Get(*s.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		set(string(kv.Key), string(kv.Value))
	}
	// 2.监听 prefix, 获得变更信息
	s.watcher(prefix, resp.Header.Revision+1, set, del)
	return nil
}

// 监听前缀中的键值对变化
func (s *ServiceDiscovery) watcher(prefix string, rev int64, set, del func(k, v string)) {
	rch := s.cli.Watch(*s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(rev))
	for watchResp := range rch {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT: // 新增&修改
				set(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE: // 删除
				del(string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}
}

//Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
