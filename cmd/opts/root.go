/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : root.go
#   Created       : 2019-01-30 10:05:09
#   Describe      :
#
# ====================================================*/
package opts

import (
	"fmt"
	"time"

	"github.com/hiruok/etcd-cli/pkg/store"
	"github.com/hiruok/etcd-cli/pkg/store/etcd"
)

// Root 获取etcd的客户端
type Root struct {
	s store.Store
}

// NewRoot 根据输入参数获取etcd客户端
func NewRoot(host string, port int32) (*Root, error) {
	e, err := etcdstore.New([]string{fmt.Sprintf("%s:%d", host, port)}, &store.Config{
		ConnectionTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	r := &Root{
		s: e,
	}
	return r, nil
}

// Close 关闭连接
func (r *Root) Close() error {
	if r.s != nil {
		return r.s.Close()
	}
	return nil
}
