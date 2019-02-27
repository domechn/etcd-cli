/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : save.go
#   Created       : 2019-01-30 10:04:20
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hiruok/etcd-cli/pkg/utils"
)

// Upload 将本地文件以二进制的形式保存到etcd
func (r *Root) Upload(etcdP, fileP string) error {
	if err := r.upload(etcdP, fileP); err != nil {
		return fmt.Errorf("upload: %s", err.Error())
	}
	return nil
}

func (r *Root) upload(etcdP, fileP string) error {
	etcdP = pathHandler(etcdP)
	var b []byte
	var err error
	if b, err = readFile(fileP); err != nil {
		return err
	}
	if err = validate(b, fileP); err != nil {
		return err
	}
	fInfo, _ := os.Stat(fileP)
	fileName := fInfo.Name()
	etcdP = etcdP + fileName
	return r.s.Put(context.Background(), utils.Normalize(etcdP), b, nil)
}
