/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : download.go
#   Created       : 2019-01-30 10:04:31
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/domgoer/etcd-cli/pkg/store"
)

// Download 将etcd中的数据下载到本地文件夹
func (r *Root) Download(etcdP, saveP string) error {
	if err := r.download(etcdP, saveP); err != nil {
		return fmt.Errorf("download: %s", err.Error())
	}
	return nil
}

func (r *Root) download(etcdP, saveP string) error {
	path := strings.TrimSuffix(pathHandler(etcdP), "/")
	k, err := r.s.Get(context.Background(), path)
	if err != nil {
		if err == store.ErrKeyNotExsit {
			return fmt.Errorf("file %s is not exsit", path)
		}
		return err
	}
	saveP = strings.TrimSuffix(saveP, "/")
	f := fmt.Sprintf("%s/%s", saveP, fileName(etcdP))
	return ioutil.WriteFile(f, k.Value, 0666)
}
