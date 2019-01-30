/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : pwd.go
#   Created       : 2019-01-30 13:14:01
#   Describe      :
#
# ====================================================*/
package opts

import "fmt"

// Pwd 返回当前路径
func (r *Root) Pwd() error {
	fmt.Println(PWD)
	return nil
}
