/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : cd.go
#   Created       : 2019-01-30 11:06:02
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"context"
	"fmt"
)

// ChangeDir 更换工作目录
func (r *Root) ChangeDir(path string) error {
	if err := r.changeDir(path); err != nil {
		return fmt.Errorf("cd: %s", err.Error())
	}
	return nil
}

func (r *Root) changeDir(path string) error {
	path = pathHandler(path)
	list, err := r.s.List(context.Background(), path)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("dir %s is not exsit or path is not dir", path)
	}
	PWD = path
	return nil
}
