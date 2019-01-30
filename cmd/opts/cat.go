/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : cat.go
#   Created       : 2019-01-30 10:27:15
#   Describe      :
#
# ====================================================*/
package opts

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hiruok/etcd-cli/pkg/store"
)

// Cat 获取文件内容
func (r *Root) Cat(path string) error {
	if err := r.cat(path); err != nil {
		return fmt.Errorf("cat: %s", err.Error())
	}
	return nil
}

func (r *Root) cat(path string) error {
	path = strings.TrimSuffix(pathHandler(path), "/")
	kv, err := r.s.Get(context.Background(), path)
	if err != nil {
		if err == store.ErrKeyNotExsit {
			return fmt.Errorf("file %s is not exsit or it is a dir", path)
		}
		return err
	}
	bio := bufio.NewReader(bytes.NewReader(kv.Value))
	var lineNum = 1
	for {
		line, err := bio.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("read file error: %v", err)
			}
		}
		fmt.Printf("%3d|\t%s", lineNum, line)
		lineNum++
	}
	return nil
}
