/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : copy.go
#   Created       : 2019-01-30 16:08:09
#   Describe      :
#
# ====================================================*/
package cmd

import "fmt"

// Copy 拷贝文件
func (r *Root) Copy(dist, src string) error {
	if err := r.copy(dist, src); err != nil {
		return fmt.Errorf("cp: %s", err.Error())
	}
	return nil
}

func (r *Root) copy(dist, src string) error {
	if pathHandler(dist) == pathHandler(src) {
		return fmt.Errorf("%s and %s are identical (not copied)", dist, src)
	}
	return r.move(dist, src, false)
}
