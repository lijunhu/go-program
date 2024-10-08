package test

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

var (
	endpoints  = []string{"10.160.84.104:2379", "10.160.84.105:2379", "10.160.84.106:2379"}
	key        = "/abc/1213212"
	revision   = 25769175
	c          = context.Background()
	cancel     = func() {}
	config     clientv3.Config
	etcdClient = new(clientv3.Client)
)

func TestRevisionGet(t *testing.T) {
	c, cancel = context.WithCancel(c)
	config = clientv3.Config{
		Endpoints:            endpoints,
		AutoSyncInterval:     300000000,
		DialTimeout:          300000000,
		DialKeepAliveTime:    300000000,
		DialKeepAliveTimeout: 300000000,
		MaxCallSendMsgSize:   10 * 1024 * 1024,
		MaxCallRecvMsgSize:   100 * 1024 * 1024,
		Context:              c,
	}
	defer cancel()

	if etcdClient, err = clientv3.New(config); err != nil {
		t.Logf("connect etcd failed.err:%s", err)
	}

	var (
		response = new(clientv3.GetResponse)
		ops      []clientv3.OpOption
		a        []clientv3.OpOption
	)

	ops = append(ops, a...)
	//ops = append(ops, clientv3.WithKeysOnly())
	//ops = append(ops, clientv3.WithPrefix())
	//ops = append(ops,clientv3.WithRev(int64(revision)))
	//ops = append(ops, clientv3.WithRev(int64(revision)))
	if response, err = etcdClient.Get(c, key, ops...); err != nil {
		t.Logf("get key: %s revision:%d failed.1err:%s", key, revision, err)
		return
	}
	r := new(clientv3.TxnResponse)
	r, err = etcdClient.Txn(c).If().Then().Else().Commit()
	if r.Succeeded {
		t.Log(1111111111111)
	}

	if len(response.Kvs) <= 0 {
		t.Logf("get key: %s revision:%d failed.2err:%s", key, revision, err)
		return
	}

	t.Log(response.Kvs[0].CreateRevision)
	t.Logf("clusterId :%x----memberId:%x", response.Header.ClusterId, response.Header.MemberId)
	var bytes []byte
	for _, kvs := range response.Kvs {
		bytes = append(bytes, kvs.Value...)
	}
	t.Log(string(bytes))
}
