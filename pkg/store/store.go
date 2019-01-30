/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : store.go
#   Created       : 2019-01-29 11:46:46
#   Describe      :
#
# ====================================================*/
package store

import (
	"context"
	"fmt"
	"time"
)

// OptType 观察到的值的变化原因
type OptType int

var (
	// ErrKeyNotExsit key值不存在
	ErrKeyNotExsit = fmt.Errorf("key not exsit")
)

const (
	// OptionTypeNew event type new
	OptionTypeNew = OptType(0)
	// OptionTypeUpdate event type update
	OptionTypeUpdate = OptType(1)
	// OptionTypeDelete event type delete
	OptionTypeDelete = OptType(2)
)

// WriteOptions put操作时的属性
type WriteOptions struct {
	IsDir     bool
	TTL       time.Duration
	KeepAlive bool
}

// LockOptions 锁的最长维持时间
type LockOptions struct {
	TTL time.Duration // Optional, expiration ttl associated with the lock
}

// KVPair 分装查询的结果
type KVPair struct {
	Key       string
	Value     []byte
	LastIndex int64
}

// WatchRes 分装watch时变化的值
type WatchRes struct {
	KV   KVPair
	Type OptType
}

// Locker 分布式锁
type Locker interface {
	// Lock 阻塞锁
	Lock()
	// Unlock 释放阻塞锁
	Unlock()
}

// NonBlockLocker 分布式非阻塞锁
type NonBlockLocker interface {
	// NonBlockLock 非阻塞式加锁 , 如果无法获取锁直接返回false
	NonBlockLock() bool
	// UnNonBlockLock 非阻塞锁解锁
	UnNonBlockLock()
}

// Store 实现存储功能
type Store interface {
	// Put 将该key的value进行修改，如果key不存在就创建
	Put(ctx context.Context, key string, value []byte, options *WriteOptions) error

	// Get 按key查询value
	Get(ctx context.Context, key string) (*KVPair, error)

	// Delete 在存储中删除该键值
	Delete(ctx context.Context, key string) error

	// Verify 查询存储中是否有该key
	Exists(ctx context.Context, key string) (bool, error)

	// Watch 观察一个key中的value变化
	Watch(ctx context.Context, key string, stopCh <-chan struct{}) (<-chan *WatchRes, error)

	// WatchTree 观察该文件夹下所有值的变化，每当有key变化时返回变化的值
	WatchTree(ctx context.Context, directory string, stopCh <-chan struct{}) (<-chan *WatchRes, error)

	// NewLock 创建一个分布式锁，但并没有锁住，如果需要加锁请调用.Lock()方法
	NewLock(key string, options *LockOptions) (Locker, error)

	// NewNonBlockLocker 返回一个非阻塞分布式锁
	NewNonBlockLocker(key string, options *LockOptions) NonBlockLocker

	// List 前缀为该值的所有kv
	List(ctx context.Context, directory string) ([]*KVPair, error)

	// DeleteTree 删除前缀为该值所有kv
	DeleteTree(ctx context.Context, directory string) error

	// AtomicPut 比较当前值和预期值，如果相同则修改
	AtomicPut(ctx context.Context, key string, value []byte, previous *KVPair, options *WriteOptions) (bool, *KVPair, error)

	// AtomicDelete 比较当前值和预期值，如果相同则删除
	AtomicDelete(ctx context.Context, key string, previous *KVPair) (bool, *KVPair, error)

	// Close 关闭连接
	Close() error
}
