/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : get.go
#   Created       : 2019-01-29 16:26:41
#   Describe      :
#
# ====================================================*/
package client

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// GetConfByEnv 根据环境变量 RUN_ENV 获取对应的配置信息,
func (c *Client) GetConfByEnv(ctx context.Context, key string) ([]byte, error) {
	env := os.Getenv(ENV)
	kv, err := c.s.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	c.fileType, err = fileType(key)
	if err != nil {
		return nil, err
	}
	return c.getConfByEnv(ctx, kv.Value, env)
}

func (c *Client) getConfByEnv(ctx context.Context, data []byte, env string) ([]byte, error) {
	var mmap map[string]interface{}
	var err error
	if mmap, err = c.unmarshal(data); err != nil {
		return nil, err
	}
	resMap, ok := mmap[env]
	if !ok {
		return nil, fmt.Errorf("not found conf's env is %s", env)
	}

	return json.Marshal(resMap)
}

func (c *Client) unmarshal(data []byte) (map[string]interface{}, error) {
	var mmap = make(map[string]interface{})
	var marshalFunc func([]byte, interface{}) error
	switch c.fileType {
	case "yaml", "yml":
		marshalFunc = yaml.Unmarshal
	case "json":
		marshalFunc = json.Unmarshal
	default:
		return nil, fmt.Errorf("unsupport type")
	}

	var err error
	if marshalFunc != nil {
		err = marshalFunc(data, &mmap)
	}
	return mmap, err
}

func fileType(path string) (string, error) {
	ss := strings.Split(path, ".")
	if len(ss) < 2 {
		return "", fmt.Errorf("key is invalid")
	}
	fileType := ss[len(ss)-1]
	return fileType, nil
}
