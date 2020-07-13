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

// WriteOptions put操作时的属性
type WriteOptions struct {
	IsDir     bool
	TTL       time.Duration
	KeepAlive bool
}

// KVPair 分装查询的结果
type KVPair struct {
	Key       string
	Value     []byte
	LastIndex int64
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
