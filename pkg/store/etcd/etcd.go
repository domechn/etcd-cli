/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : etcd.go
#   Created       : 2019-01-29 11:48:14
#   Describe      :
#
# ====================================================*/
package etcdstore

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/domgoer/etcd-cli/pkg/store"
)

// Etcd implements Store
type Etcd struct {
	client clientv3.Client
}

// etcdLock implements Locker
type etcdLock struct {
	client clientv3.Client
	key    string
	ttl    time.Duration
	locker sync.Locker
}

type nonBlockEtcdLock struct {
	client clientv3.Client
	key    string
	ttl    time.Duration
}

const (
	periodicSync      = 5 * time.Minute  // 自动同步间隔
	defaultLockTTL    = 20 * time.Second // 锁的默认超时时间
	defaultUpdateTime = 5 * time.Second
	// DefaultLockKey 锁的默认key
	DefaultLockKey = "mutex/lock/key"
)

// New 通过给定的地址列表和tls配置，创建一个新的etcd客户端
func New(addrs []string, options *store.Config) (*Etcd, error) {
	s := &Etcd{}

	var (
		entries []string
		err     error
	)

	entries = store.CreateEndpoints(addrs, "http")

	cfg := &clientv3.Config{
		Endpoints:   entries,
		DialTimeout: 10 * time.Second,
		// 不自动同步
		// AutoSyncInterval: periodicSync,
	}
	// 设置 options
	if options != nil {
		if options.TLS != nil {
			setTLS(cfg, options.TLS, addrs)
		}
		if options.ConnectionTimeout != 0 {
			setTimeout(cfg, options.ConnectionTimeout)
		}
		if options.Username != "" {
			setCredentials(cfg, options.Username, options.Password)
		}
	}

	c, err := clientv3.New(*cfg)
	if err != nil {
		return nil, err
	}
	s.client = *c

	return s, nil

}

// SetTLS 设置证书路径
func setTLS(cfg *clientv3.Config, tls *tls.Config, addrs []string) {
	entries := store.CreateEndpoints(addrs, "https")
	cfg.Endpoints = entries

	cfg.TLS = tls
}

// setTimeout 设置etcd连接超时时间
func setTimeout(cfg *clientv3.Config, time time.Duration) {
	cfg.DialTimeout = time
}

// setCredentials 用于设置https认证
func setCredentials(cfg *clientv3.Config, username, password string) {
	cfg.Username = username
	cfg.Password = password
}

// Normalize 转换这个值用于etcd
func (s *Etcd) normalize(key string) string {
	key = store.Normalize(key)
	return strings.TrimPrefix(key, "/")
}

// NormalizeDir 如果值中是目录，就检查参数是否符合目录格式
// 如果不符合就转化
func (s *Etcd) normalizeDir(key string) string {
	tk := s.normalize(key)
	if !strings.HasSuffix(tk, "/") {
		tk += "/"
	}
	return tk
}

// Put 更新该key的value,如果超时时间<=0，则没有超时时间
func (s *Etcd) Put(ctx context.Context, key string, value []byte, opts *store.WriteOptions) error {
	var op []clientv3.OpOption
	keys := s.normalize(key)
	if opts != nil {
		if opts.TTL > 0 {
			lease, err := s.client.Grant(ctx, int64(opts.TTL.Seconds()))
			if err != nil {
				return err
			}
			op = append(op, clientv3.WithLease(lease.ID))
			if opts.KeepAlive {
				// 如果keepalive没有挂，那么key就一直存在，如果keepalie挂了，超过ttl，key就消失了
				_, err := s.client.KeepAlive(context.Background(), lease.ID)
				if err != nil {
					return err
				}
			}
		}
		if opts.IsDir {
			keys += "/"
		}
	}
	_, err := s.client.Put(ctx, keys, string(value), op...)
	return err
}

// Get 获取一个key的值和它最后一次修改的版本号
// 这个版本号用来Atomic操作时的cas
func (s *Etcd) Get(ctx context.Context, key string) (*store.KVPair, error) {
	result, err := s.client.Get(ctx, s.normalize(key))
	if err != nil {
		return nil, err
	}
	kvs := result.Kvs
	if len(result.Kvs) == 0 {
		return nil, store.ErrKeyNotExsit
	}
	kv := kvs[0]
	resKV := &store.KVPair{
		Key:       key,
		Value:     kv.Value,
		LastIndex: kv.ModRevision,
	}
	return resKV, nil
}

// Delete 删除这个key和它的值
func (s *Etcd) Delete(ctx context.Context, key string) error {
	_, err := s.client.Delete(ctx, s.normalize(key))
	return err
}

// Exists 查看该key是否存在在etcd中
func (s *Etcd) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.Get(ctx, key)
	if err != nil {
		if err != store.ErrKeyNotExsit {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// Watch 观察一个key的变化
// 它返回一个chan可以用来接收该key的value的变化
// 向stopch中传值可以终止watch
func (s *Etcd) Watch(ctx context.Context, key string, stopCh <-chan struct{}) (<-chan *store.WatchRes, error) {
	rch := s.client.Watch(ctx, s.normalize(key))
	watchCh := make(chan *store.WatchRes)
	go func() {
		defer close(watchCh)
		for {
			select {
			case <-stopCh:
				return
			case wresp := <-rch:
				for _, ev := range wresp.Events {
					var optType store.OptType
					// 判断操作类型
					switch ev.Type {
					case mvccpb.DELETE:
						optType = store.OptionTypeDelete
					case mvccpb.PUT:
						if ev.IsCreate() {
							optType = store.OptionTypeNew
						} else if ev.IsModify() {
							optType = store.OptionTypeUpdate
						}
					}
					watchCh <- &store.WatchRes{
						KV: store.KVPair{
							Key:       string(ev.Kv.Key),
							Value:     ev.Kv.Value,
							LastIndex: ev.Kv.ModRevision,
						},
						Type: optType,
					}
				}
			}
		}
	}()

	return watchCh, nil
}

// WatchTree 监视“目录”上的更改
// 它返回一个chan可以用来接收directory下kv的变化
// 向stopch中传值可以终止watchtree
func (s *Etcd) WatchTree(ctx context.Context, directory string, stopCh <-chan struct{}) (<-chan *store.WatchRes, error) {
	rch := s.client.Watch(ctx, s.normalizeDir(directory), clientv3.WithPrefix())
	watchCh := make(chan *store.WatchRes)

	go func() {
		defer close(watchCh)

		for {
			// Check if the watch was stopped by the caller
			select {
			case <-stopCh:
				return
			case wresp := <-rch:
				for _, ev := range wresp.Events {
					var optType store.OptType
					// 判断操作类型
					switch ev.Type {
					case mvccpb.DELETE:
						optType = store.OptionTypeDelete
					case mvccpb.PUT:
						if ev.IsCreate() {
							optType = store.OptionTypeNew
						} else if ev.IsModify() {
							optType = store.OptionTypeUpdate
						}
					}
					watchCh <- &store.WatchRes{
						KV: store.KVPair{
							Key:       string(ev.Kv.Key),
							Value:     ev.Kv.Value,
							LastIndex: ev.Kv.ModRevision,
						},
						Type: optType,
					}
				}
			}

		}
	}()

	return watchCh, nil

}

// NewLock 返回一个所得结构体，可以用来锁住一个key
func (s *Etcd) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	ttl := defaultLockTTL

	// 设置锁的过期时间
	if options != nil {
		if options.TTL != 0 {
			ttl = options.TTL
		}
	}

	// 创建一个锁的对象
	l := &etcdLock{
		client: s.client,
		key:    s.normalize(key),
		ttl:    ttl,
	}

	session, err := concurrency.NewSession(&l.client, concurrency.WithTTL(int(l.ttl.Seconds())))
	if err != nil {
		return nil, err

	}
	locker := concurrency.NewLocker(session, l.key)
	l.locker = locker

	return l, nil
}

// List 列出该directory下的所有kv
func (s *Etcd) List(ctx context.Context, directory string) ([]*store.KVPair, error) {
	var resKV []*store.KVPair
	result, err := s.client.Get(ctx, s.normalize(directory), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range result.Kvs {
		res := &store.KVPair{
			Key:       string(v.Key),
			Value:     v.Value,
			LastIndex: v.ModRevision,
		}
		resKV = append(resKV, res)
	}
	return resKV, nil
}

// DeleteTree 产出directory下所有的键值
func (s *Etcd) DeleteTree(ctx context.Context, directory string) error {
	_, err := s.client.Delete(ctx, s.normalizeDir(directory), clientv3.WithPrefix())
	return err
}

// AtomicPut 更新一个值在“键”，如果被键同时修改，则抛出错误
// 如果操作失败，它将返回(false, previous, error (if err != nil))
// 如果操作成功，它将返回(true, nil, nil)
func (s *Etcd) AtomicPut(ctx context.Context, key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) {

	var op []clientv3.OpOption
	keys := s.normalize(key)
	if opts != nil {
		if opts.TTL > 0 {
			lease, err := s.client.Grant(ctx, int64(opts.TTL.Seconds()))
			if err != nil {
				return false, previous, err
			}
			op = append(op, clientv3.WithLease(lease.ID))
		}
		if opts.IsDir {
			if !strings.HasSuffix(keys, "/") {
				keys += "/"
			}
		}
	}

	resp, err := s.client.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", int64(previous.LastIndex))).
		Then(clientv3.OpPut(key, string(value), op...)).
		Else(clientv3.OpGet(key)).
		Commit()
	if err != nil {
		return false, previous, err
	}

	if !resp.Succeeded {
		res := resp.Responses[0]
		getRes := res.GetResponseRange()
		if len(getRes.Kvs) == 0 {
			return false, previous, nil
		}
		kv := getRes.Kvs[0]
		pair := &store.KVPair{
			Key:       key,
			Value:     kv.Value,
			LastIndex: kv.ModRevision,
		}
		return false, pair, nil
	}

	return true, nil, nil
}

// AtomicDelete 删除该键的值，如果这个键被同时修改，抛出一个错误
func (s *Etcd) AtomicDelete(ctx context.Context, key string, previous *store.KVPair) (bool, *store.KVPair, error) {

	resp, err := s.client.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", int64(previous.LastIndex))).
		Then(clientv3.OpDelete(key)).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return false, previous, err
	}

	if !resp.Succeeded {
		res := resp.Responses[0]
		getRes := res.GetResponseRange()

		if len(getRes.Kvs) == 0 {
			return true, nil, nil
		}

		kv := getRes.Kvs[0]
		pair := &store.KVPair{
			Key:       key,
			Value:     kv.Value,
			LastIndex: kv.ModRevision,
		}
		return false, pair, nil
	}

	return true, nil, nil

}

// Close 关闭客户端连接
func (s *Etcd) Close() error {
	return s.client.Close()
}

// Lock 用etc，加锁这个key
func (l *etcdLock) Lock() {
	l.locker.Lock()
}

// Unlock 解锁这个key.
func (l *etcdLock) Unlock() {
	l.locker.Unlock()
}

// NewNonBlockLocker 返回一个非阻塞分布式锁
func (s *Etcd) NewNonBlockLocker(key string, options *store.LockOptions) store.NonBlockLocker {
	ttl := defaultLockTTL

	// 设置锁的过期时间
	if options != nil {
		if options.TTL != 0 {
			ttl = options.TTL
		}
	}

	// 创建一个锁的对象
	l := &nonBlockEtcdLock{
		client: s.client,
		key:    s.normalize(key),
		ttl:    ttl,
	}

	return l
}

func (l *nonBlockEtcdLock) NonBlockLock() bool {
	ctx := context.Background()
	var op clientv3.OpOption
	lease, err := l.client.Grant(ctx, int64(l.ttl.Seconds()))
	if err != nil {
		return false
	}
	op = clientv3.WithLease(lease.ID)

	resp, err := l.client.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(l.key), "=", 0)).
		Then(clientv3.OpPut(l.key, "", op)).
		Else(clientv3.OpGet(l.key)).
		Commit()
	if err != nil || !resp.Succeeded {
		return false
	}

	return true
}

func (l *nonBlockEtcdLock) UnNonBlockLock() {
	l.client.Delete(context.Background(), l.key)
}
