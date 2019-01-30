/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : client.go
#   Created       : 2019-01-29 15:56:31
#   Describe      :
#
# ====================================================*/
package client

import (
	"strings"
	"time"

	"github.com/hiruok/etcd-cli/pkg/tls"

	"github.com/hiruok/etcd-cli/pkg/store"
	"github.com/hiruok/etcd-cli/pkg/store/etcd"
)

const (
	defaultTimeout = 10 * time.Second
)

// Client etcd客户端
type Client struct {
	s        store.Store
	fileType string
}

// Config 用于配置连接tls服务端的信息
type Config struct {
	// ca 文件的路径
	CA string
	// cert 文件的路径
	Cert string
	// key 文件的路径
	Key string

	// 连接用的账号密码
	Username string
	Password string
}

// NewClient 返回一个新的客户端对象，如果连接失败返回错误，默认连接超时时间10s
// 如果有多个地址，用";"分割，如"127.0.0.1:2379;127.0.0.1:2378;..."
func NewClient(addrs string) (*Client, error) {
	e, err := etcdstore.New(strings.Split(addrs, ";"), &store.Config{
		ConnectionTimeout: defaultTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		s: e,
	}, nil
}

// NewTLSClient 返回一个连接到tls的etcd服务的新的客户端对象，
// cfg中配置连接用的证书的位置,其他参数和NewClient中一致
func NewTLSClient(addrs string, cfg Config) (*Client, error) {
	tlsC, err := tlsclient.Config(tlsclient.Options{
		CAFile:   cfg.CA,
		CertFile: cfg.Cert,
		KeyFile:  cfg.Key,
	})
	if err != nil {
		return nil, err
	}
	e, err := etcdstore.New(strings.Split(addrs, ";"), &store.Config{
		ConnectionTimeout: defaultTimeout,
		TLS:               tlsC,
		Username:          cfg.Username,
		Password:          cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		s: e,
	}, nil
}

// Close 关闭和etcd的连接
func (c *Client) Close() error {
	return c.s.Close()
}
