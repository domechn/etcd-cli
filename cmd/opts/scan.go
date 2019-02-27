/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : scan.go
#   Created       : 2019-01-30 10:14:21
#   Describe      :
#
# ====================================================*/
package opts

import (
	"fmt"
)

// DoScan 扫描输入内容，并调用对应的函数
func (r *Root) DoScan(bs []string) error {
	return r.doScan(bs)
}

func (r *Root) doScan(bs []string) error {
	defer func() {
		if err := recover(); err != nil {
			if len(bs) != 0 {
				fmt.Println(ErrInvalidParamNum.Error())
			}
		}
	}()
	var err error
	switch bs[0] {
	case "upload":
		err = r.Upload(bs[1], bs[2])
	case "download":
		err = r.Download(bs[1], bs[2])
	case "ls":
		var p string
		if len(bs) > 1 {
			p = bs[1]
		}
		err = r.List(p)
	case "cd":
		err = r.ChangeDir(bs[1])
	case "pwd":
		err = r.Pwd()
	case "mkdir":
		err = r.Mkdir(bs[1])
	case "cat":
		err = r.Cat(bs[1])
	case "vim":
		err = r.Vim(bs[1])
	case "rm":
		var params []string
		for i := 1; i < len(bs); i++ {
			params = append(params, bs[i])
		}
		err = r.Remove(params...)
	case "touch":
		err = r.Touch(bs[1])
	case "mv":
		err = r.Move(bs[1], bs[2])
	case "cp":
		err = r.Copy(bs[1], bs[2])
	default:
		err = fmt.Errorf("unsupport")
	}
	return err
}
