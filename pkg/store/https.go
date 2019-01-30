/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : https.go
#   Created       : 2019-01-29 11:47:00
#   Describe      :
#
# ====================================================*/
package store

import (
	"crypto/tls"
	"time"
)

// Config contains the options for a storage client
type Config struct {
	TLS               *tls.Config
	ConnectionTimeout time.Duration
	Username          string
	Password          string
}
