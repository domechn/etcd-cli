/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : move.go
#   Created       : 2019-01-30 15:15:24
#   Describe      :
#
# ====================================================*/
package opts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hiruok/etcd-cli/pkg/store"
)

// Move 源地址，目标地址
func (r *Root) Move(dist, src string) error {
	if err := r.move(dist, src, true); err != nil {
		return fmt.Errorf("mv: %s", err.Error())
	}
	return nil
}

func (r *Root) move(dist, src string, delete bool) error {
	var distIsDir, srcIsDir bool

	if strings.HasSuffix(dist, "/") {
		distIsDir = true
	}

	if strings.HasSuffix(src, "/") {
		srcIsDir = true
	}
	dist = pathHandler(dist)
	src = pathHandler(src)
	if distIsDir {
		if !srcIsDir {
			return fmt.Errorf("you can't mv a dir to a file")
		}
		newKey := src

		kvs, err := r.s.List(context.Background(), dist)
		if err != nil {
			return err
		}

		newKVs, err := r.s.List(context.Background(), newKey)
		if err != nil {
			return err
		}
		if len(newKVs) != 0 {
			newKey = newKey + fileName(strings.TrimSuffix(dist, "/")) + "/"
		}

		for _, kv := range kvs {
			k := strings.Replace(kv.Key, dist, newKey, 1)
			err = r.s.Put(context.Background(), k, kv.Value, nil)
			if err != nil {
				return err
			}
		}
		if delete {
			err = r.s.DeleteTree(context.Background(), dist)
		}
		return err
	}

	dist = strings.TrimSuffix(dist, "/")
	fileName := fileName(dist)

	kv, err := r.s.Get(context.Background(), dist)
	if srcIsDir {
		if err != nil {
			return nil
		}
		src = src + fileName
		if err = r.s.Put(context.Background(), src, kv.Value, nil); err != nil {
			return err
		}
		if delete {
			err = r.s.Delete(context.Background(), dist)
		}
		return err
	}

	src = strings.TrimSuffix(src, "/")
	_, err = r.s.Get(context.Background(), src)
	if err != store.ErrKeyNotExsit {
		return fmt.Errorf("file %s is exsit", src)
	}

	err = r.s.Put(context.Background(), src, kv.Value, nil)
	if err != nil {
		return err
	}
	if delete {
		err = r.s.Delete(context.Background(), dist)
	}
	return err
}
