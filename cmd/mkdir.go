/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : mkdir.go
#   Created       : 2019-01-30 13:16:59
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"context"
	"fmt"

	"github.com/hiruok/etcd-cli/pkg/store"
)

// Mkdir 创建文件夹
func (r *Root) Mkdir(path string) error {
	if err := r.mkdir(path); err != nil {
		return fmt.Errorf("mkdir: %s", err.Error())
	}
	return nil
}

func (r *Root) mkdir(path string) error {
	path = pathHandler(path)
	return r.s.Put(context.Background(), path, []byte(""), &store.WriteOptions{
		IsDir: true,
	})
}
