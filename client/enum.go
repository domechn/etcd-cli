/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : enum.go
#   Created       : 2019-01-29 16:30:34
#   Describe      :
#
# ====================================================*/
package client

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (

	// ENV 获取环境变量值
	ENV = "RUN_ENV"
	// ENVQA  qa环境变量值
	ENVQA = "qa"
	// ENVUAT uat环境变量值
	ENVUAT = "uat"
	// ENVPRO 生成环境变量值
	ENVPRO = "pro"
)
