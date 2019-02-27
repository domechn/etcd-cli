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
package opts

import (
	"context"
	"fmt"
	"strings"
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
	list, err := r.s.List(context.Background(), fp)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("dir %s is not exsit", fp)
	}

	var set = make(map[string]int8, len(list))
	var idx int
	for _, l := range list {
		if idx != 0 && idx%8 == 0 {
			changeLine()
		}
		k := pathResHandler(fp, l.Key)
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
	return nil
}

func pathResHandler(pwd, p string) string {
	p = strings.TrimPrefix(p, pwd)
	ps := strings.Split(p, "/")
	if len(ps) > 1 {
		return ps[0] + "/"
	}
	return ps[0]
}
