package discovery_etcd

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/registry"
)

var _ registry.Watcher = (*Watcher)(nil)

type Watcher struct {
	key         string
	ctx         context.Context
	cancel      context.CancelFunc
	Client      *clientv3.Client
	watchChan   clientv3.WatchChan
	watcher     clientv3.Watcher
	kv          clientv3.KV
	first       bool
	ServiceName string
}

func newWatcher(ctx context.Context, key, name string, client *clientv3.Client) (*Watcher, error) {
	w := &Watcher{
		key:         key,
		Client:      client,
		watcher:     clientv3.NewWatcher(client),
		kv:          clientv3.NewKV(client),
		first:       true,
		ServiceName: name,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0), clientv3.WithKeysOnly())
	err := w.watcher.RequestProgress(w.ctx)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Watcher) Next() ([]*registry.ServiceInstance, error) {
	if w.first {
		item, err := w.getInstance()
		w.first = false
		return item, err
	}

	select {
	default:
		return nil, nil
	case <-w.ctx.Done():
		defer w.Stop()
		return nil, w.ctx.Err()
	case watchResp, ok := <-w.watchChan:
		if !ok || watchResp.Err() != nil {
			time.Sleep(time.Second)
			err := w.reWatch()
			if err != nil {
				return nil, err
			}
		}
		i, err := w.getInstance()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("i: %v\n", i)

		return i, err
	}
}

func (w *Watcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}

func (w *Watcher) getInstance() ([]*registry.ServiceInstance, error) {
	resp, err := w.kv.Get(w.ctx, w.key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != w.ServiceName {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

func (w *Watcher) reWatch() error {
	w.watcher.Close()
	w.watcher = clientv3.NewWatcher(w.Client)
	w.watchChan = w.watcher.Watch(w.ctx, w.key, clientv3.WithPrefix(), clientv3.WithRev(0), clientv3.WithKeysOnly())
	return w.watcher.RequestProgress(w.ctx)
}
