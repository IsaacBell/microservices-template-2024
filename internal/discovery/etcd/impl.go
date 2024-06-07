package discovery_etcd

import (
	"context"
	zapWrapper "core/internal/logs"
	"core/internal/util"
	"fmt"
	"log"
	"time"

	zap "go.uber.org/zap"

	registry "github.com/go-kratos/etcd/registry"
	etcdClient "go.etcd.io/etcd/client/v3"
)

var (
	etcdRegistrar *registry.Registry
	client        *etcdClient.Client
	zlogger       *zap.Logger
	watcher       *Watcher
	watchers      []*Watcher = make([]*Watcher, 0)
	endpoints     []string   = []string{"http://localhost:2379"}
)

func GetClient() *etcdClient.Client {
	return client
}

// Create and register etcd client
func Register(ctx context.Context, name string) (*etcdClient.Client, *registry.Registry) {
	zapLogger, _ := zapWrapper.NewZlog()

	etcd, err := etcdClient.New(etcdClient.Config{
		Endpoints: endpoints,
		Context:   ctx,
		Logger:    zapLogger,
	})
	if err != nil {
		util.PrintLnInColor(util.AnsiColorRed, "Failed to create etcd client")
		log.Fatalf("%v", err)
	}

	_, err = etcd.Put(ctx, "/system/last_registered_service", name)
	if err != nil {
		util.PrintLnInColor(
			util.AnsiColorRed,
			"Failed: ectd.Put(ctx, \"/system/last_registered_service\", "+name+")\n-> ",
			err)
	} else {
		util.PrintLnInColor(
			util.AnsiColorGreen,
			"etcd: PUT /system/last_registered_service", name,
		)
	}

	client = etcd
	zlogger = zapLogger
	etcdRegistrar = registry.New(etcd, registry.Context(ctx))
	return client, etcdRegistrar
}

// stores a key-value pair in etcd
func Store(ctx context.Context, etcd *etcdClient.Client, key, value string) error {
	_, err := etcd.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to store key-value pair: %v", err)
	}
	return nil
}

// retrieves the value associated with a key from etcd
func Retrieve(ctx context.Context, etcd *etcdClient.Client, key string) (string, error) {
	resp, err := etcd.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve value: %v", err)
	}
	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}
	return "", fmt.Errorf("key not found")
}

// deletes a key from etcd
func Delete(ctx context.Context, etcd *etcdClient.Client, key string) error {
	_, err := etcd.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete key: %v", err)
	}
	return nil
}

// usage:
//
//	watcher, err := discovery_etcd.NewWatcher(ctx, "/your/service/key/prefix", serviceName, discovery_etcd.GetClient())
//	defer watcher.Stop()
func NewWatcher(ctx context.Context, name string, client *etcdClient.Client) (*Watcher, error) {
	w, err := newWatcher(ctx, "/"+name+"/service", name, client)
	if err != nil {
		log.Fatalln("Failed to create service discovery watcher: ", err)
	}
	watcher = w
	watchers = append(watchers, w)
	return w, err
}

// watches for changes to a key in etcd and passes the events to a callback function
func Watch(ctx context.Context, etcd *etcdClient.Client, key string, callback func(event *etcdClient.Event)) {
	util.PrintLnInColor(util.AnsiColorMagenta, "\nWatching key: ", util.AnsiColorCyan, key)
	for {
		util.PrintLnInColor(util.AnsiColorMagenta, "\nChecking key: ", util.AnsiColorCyan, key)
		watchChan := etcd.Watch(ctx, key)
		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				util.PrintLnInColor(util.AnsiColorRed, "Error watching etcd key:", watchResp.Err())
				break
			}
			for _, event := range watchResp.Events {
				util.PrintLnInColor(util.AnsiColorGreen, "etcd event: ", util.AnsiColorCyan, event.Kv)
				callback(event)
			}
		}
		util.PrintLnInColor(util.AnsiColorYellow, "Watch channel closed. Reestablishing watch...")
		time.Sleep(time.Second)
		go Watch(ctx, etcd, key, callback)
	}
}

// continuously watch for service instances matching the watcher's target prefix
func (w *Watcher) StartDiscovery() {
	watcher = w
	util.PrintLnInColor(util.AnsiColorGray, "etcd: Starting service discovery...")

	for {
		defer time.Sleep(time.Second * 2)
		util.PrintLnInColor(util.AnsiColorGreen, "etcd: Service discovery active...")
		servicesFound := true

		instances, err := w.Next()
		if err != nil {
			util.PrintLnInColor(util.AnsiColorYellow, "etcd: No new service discovery instances: ", util.AnsiColorRed, err)
			servicesFound = false
			continue
		}

		if len(instances) == 0 {
			servicesFound = false
		}

		if !servicesFound {
			time.Sleep(time.Second)
			continue
		}

		// Process the service instances
		for _, instance := range instances {
			util.PrintLnInColor(util.AnsiColorCyan, "Service instance: %+v", util.AnsiColorGreen, instance)
			// Perform actions based on the discovered service instances
		}
	}
}
