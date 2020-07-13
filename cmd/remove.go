/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : remove.go
#   Created       : 2019-01-30 14:40:09
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/domgoer/etcd-cli/pkg/store"
)

// Remove 删除一个文件 如果是文件夹需要在路径最后加上"/""
func (r *Root) Remove(paths ...string) error {
	r.CleanCache()
	for _, p := range paths {
		if err := r.remove(p); err != nil {
			fmt.Println("remove: ", err.Error())
		}
	}
	return nil
}

func (r *Root) remove(path string) error {
	var isDir bool
	if strings.HasSuffix(path, "/") {
		isDir = true
	}

	path = pathHandler(path)
	if !isDir {
		path = strings.TrimSuffix(path, "/")
		if _, err := r.s.Get(context.Background(), path); err == store.ErrKeyNotExsit {
			return fmt.Errorf("file %s is not exsit", path)
		}
		return r.s.Delete(context.Background(), path)
	}

	list, err := r.s.List(context.Background(), path)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("dir %s is not exsit", path)
	}
	return r.s.DeleteTree(context.Background(), path)
}
