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
package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/hiruok/etcd-cli/pkg/tls"
	"github.com/patrickmn/go-cache"
	"time"

	"github.com/hiruok/etcd-cli/pkg/store"
	"github.com/hiruok/etcd-cli/pkg/store/etcd"
)

// Root 获取etcd的客户端
type Root struct {
	s store.Store

	dirCache *cache.Cache
}

type Config struct {
	Host     string
	Port     int32
	Ca       string
	Cert     string
	Key      string
	Username string
	Password string
}

// NewRoot 根据输入参数获取etcd客户端
func NewRoot(cfg Config) (*Root, error) {
	var t *tls.Config
	if cfg.Ca != "" && cfg.Cert != "" && cfg.Key != "" {
		var err error
		t, err = tlsclient.Config(tlsclient.Options{
			CAFile:   cfg.Ca,
			CertFile: cfg.Cert,
			KeyFile:  cfg.Key,
		})
		if err != nil {
			return nil, err
		}
	}
	e, err := etcdstore.New([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, &store.Config{
		TLS:               t,
		ConnectionTimeout: 10 * time.Second,
		Username:          cfg.Username,
		Password:          cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	r := &Root{
		s:        e,
		dirCache: cache.New(time.Second*3, time.Second*10),
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

func (r *Root) CleanCache() {
	r.dirCache.Flush()
}
