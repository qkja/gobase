package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/extend/etcd"
	"github.com/qkja/gobase/isc"
	"github.com/qkja/gobase/time"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Test1(t *testing.T) {
	config.LoadYamlFile("./application-test1.yaml")
	if config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			return
		}
	}

	etcdClient, _ := etcd.NewEtcdClient()

	ctx := context.Background()
	etcdClient.Put(ctx, "test", time.TimeToStringYmdHms(time.Now()))
	rsp, _ := etcdClient.Get(ctx, "test")
	etcdClient.Get(ctx, "test", func(pOp *clientv3.Op) {
		fmt.Println("信息")
		fmt.Println(isc.ToJsonString(&pOp))
	})
	fmt.Println(rsp)
}

func TestRetry(t *testing.T) {
	config.LoadYamlFile("./application-retry.yaml")
	if config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			return
		}
	}

	_, err := etcd.NewEtcdClient()
	if err != nil {
		fmt.Println("")
	} else {
		fmt.Println("链接etcd 成功")
	}
}
