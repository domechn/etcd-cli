/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : ls.go
#   Created       : 2019-01-30 10:25:21
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// List 列出当前路径或者指定路径fp下所有的文件
func (r *Root) List(fp string) error {
	if err := r.list(fp); err != nil {
		return fmt.Errorf("list: %s", err.Error())
	}
	return nil
}

func (r *Root) list(fp string) error {
	fp = pathHandler(fp)
	if names, ok := r.dirCache.Get(fp); ok {
		printFiles(fp, names.([]string))
		return nil
	}
	list, err := r.s.List(context.Background(), fp)
	var files []string

	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("dir %s is not exsit", fp)
	}
	for _, l := range list {
		files = append(files, l.Key)
	}

	printFiles(fp, files)
	r.dirCache.Set(fp, files, time.Second*2)

	return nil
}

func printFiles(fp string, files []string) {
	var set = make(map[string]int8, len(files))
	var idx int

	for _, f := range files {
		if idx != 0 && idx%8 == 0 {
			changeLine()
		}
		k := pathResHandler(fp, f)
		if k == "/" || k == "" {
			continue
		}
		if _, ok := set[k]; ok {
			continue
		}
		fmt.Printf("%s\t", k)
		set[k] = 1
		idx++
	}
	changeLine()
}

func pathResHandler(pwd, p string) string {
	p = strings.TrimPrefix(p, pwd)
	ps := strings.Split(p, "/")
	if len(ps) > 1 {
		return ps[0] + "/"
	}
	return ps[0]
}
