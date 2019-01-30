// Copyright (c) 2018, dmc (814172254@qq.com),
//
// Authors: dmc,
//
// Distribution:生成一个tls客户端用于连接https服务端.
package tlsclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"runtime"
)

// Client TLS密码组件(为客户首选的组件提供CBC密码)
var clientCipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

// Options 用来创建tls客户端的配置信息
type Options struct {
	CAFile string

	// 文件路径，不能为空
	CertFile string
	KeyFile  string

	// 如果该值为true并且CAFile存在的话，那么用于tls组件的root pool将使用CAFile中的roots
	// 否则使用系统的pool
	ExclusiveRootPools bool
	MinVersion         uint16
}

// allTLSVersions 列出所有的tls版本
var allTLSVersions = map[uint16]struct{}{
	tls.VersionSSL30: {},
	tls.VersionTLS10: {},
	tls.VersionTLS11: {},
	tls.VersionTLS12: {},
}

// ClientDefault 返回一个默认的tls配置
func ConfigDefault(ops ...func(*tls.Config)) *tls.Config {
	tlsconfig := &tls.Config{
		// Prefer TLS1.2 as the client minimum
		MinVersion:   tls.VersionTLS12,
		CipherSuites: clientCipherSuites,
	}

	for _, op := range ops {
		op(tlsconfig)
	}

	return tlsconfig
}

// certPool 返回一个X.509certPool从cafile中
func certPool(caFile string, exclusivePool bool) (*x509.CertPool, error) {
	// If we should verify the server, we need to load a trusted ca
	var (
		certPool *x509.CertPool
		err      error
	)
	if exclusivePool {
		certPool = x509.NewCertPool()
	} else {
		certPool, err = SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to read system certificates: %v", err)
		}
	}
	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate %q: %v", caFile, err)
	}
	if !certPool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("failed to append certificates from PEM file: %q", caFile)
	}
	return certPool, nil
}

// isValidMinVersion 从tls所有版本中查看提供的版本是否合法
func isValidMinVersion(version uint16) bool {
	_, ok := allTLSVersions[version]
	return ok
}

// adjustMinVersion 设置配置的版本，option.config中的MinVersion必须比tls.config中的MinVersion大，否则返回版本过低
func adjustMinVersion(options Options, config *tls.Config) error {
	if options.MinVersion > 0 {
		if !isValidMinVersion(options.MinVersion) {
			return fmt.Errorf("Invalid minimum TLS version: %x\n", options.MinVersion)
		}
		if options.MinVersion < config.MinVersion {
			return fmt.Errorf("Requested minimum TLS version is too low. Should be at-least: %x\n", config.MinVersion)
		}
		config.MinVersion = options.MinVersion
	}

	return nil
}

// getCert 返回从certfile和keyfile中读取出的Certificate
func getCert(options Options) ([]tls.Certificate, error) {
	if options.CertFile == "" && options.KeyFile == "" {
		return nil, nil
	}

	errMessage := "Could not load X509 key pair"

	tlsCert, err := tls.LoadX509KeyPair(options.CertFile, options.KeyFile)

	if err != nil {
		return nil, fmt.Errorf("%v\n%s\n", err, errMessage)
	}

	return []tls.Certificate{tlsCert}, nil
}

// Client 返回一个可以被客户端使用tls配置
func Config(options Options) (*tls.Config, error) {
	tlsConfig := ConfigDefault()
	if options.CAFile != "" {
		CAs, err := certPool(options.CAFile, options.ExclusiveRootPools)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = CAs
	}

	tlsCerts, err := getCert(options)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = tlsCerts

	if err := adjustMinVersion(options, tlsConfig); err != nil {
		return nil, err
	}

	return tlsConfig, nil
}

// SystemCertPool 返回系统的证书池
func SystemCertPool() (*x509.CertPool, error) {
	certpool, err := x509.SystemCertPool()
	if err != nil && runtime.GOOS == "windows" {
		return x509.NewCertPool(), nil
	}
	return certpool, err
}
