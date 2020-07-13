/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : touch.go
#   Created       : 2019-01-30 15:23:43
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

// Touch 创建一个文件
func (r *Root) Touch(path string) error {
	r.CleanCache()
	if err := r.touch(path); err != nil {
		return fmt.Errorf("touch: %s", err.Error())
	}
	return nil
}

func (r *Root) touch(path string) error {
	path = strings.TrimSuffix(pathHandler(path), "/")
	if _, err := r.s.Get(context.Background(), path); err != store.ErrKeyNotExsit {
		return fmt.Errorf("file is exsit")
	}

	return r.s.Put(context.Background(), path, []byte(""), nil)
}
