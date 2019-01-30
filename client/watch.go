/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : watch.go
#   Created       : 2019-01-29 16:41:41
#   Describe      :
#
# ====================================================*/
package client

import (
	"context"
	"fmt"
	"os"
)

// WatchConfByEnv 观察env对应的配置信息是否发生变化
func (c *Client) WatchConfByEnv(ctx context.Context, key string, stopCh <-chan struct{}) (<-chan []byte, error) {
	res, err := c.s.Watch(ctx, key, stopCh)
	if err != nil {
		return nil, err
	}

	var resByte = make(chan []byte)

	go func() {
		for {
			w := <-res
			fmt.Println(string(w.KV.Value))
			b, err := c.getConfByEnv(context.Background(), w.KV.Value, os.Getenv(ENV))
			if err == nil {
				resByte <- b
			} else {
				fmt.Println(err)
			}
		}
	}()

	return resByte, nil
}
