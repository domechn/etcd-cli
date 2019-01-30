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
func (r *Root) DoScan(bs [][]byte) error {
	return r.doScan(bs)
}

func (r *Root) doScan(bs [][]byte) error {
	defer func() {
		if err := recover(); err != nil {
			if len(bs) != 0 {
				fmt.Println(ErrInvalidParamNum.Error())
			}
		}
	}()
	var err error
	switch string(bs[0]) {
	case "upload":
		err = r.Upload(string(bs[1]), string(bs[2]))
	case "download":
		err = r.Download(string(bs[1]), string(bs[2]))
	case "ls":
		var p string
		if len(bs) > 1 {
			p = string(bs[1])
		}
		err = r.List(p)
	case "cd":
		err = r.ChangeDir(string(bs[1]))
	case "pwd":
		err = r.Pwd()
	case "mkdir":
		err = r.Mkdir(string(bs[1]))
	case "cat":
		err = r.Cat(string(bs[1]))
	case "vim":
		err = r.Vim(string(bs[1]))
	case "rm":
		var params []string
		for i := 1; i < len(bs); i++ {
			params = append(params, string(bs[i]))
		}
		err = r.Remove(params...)
	case "touch":
		err = r.Touch(string(bs[1]))
	case "mv":
		err = r.Move(string(bs[1]), string(bs[2]))
	case "cp":
		err = r.Copy(string(bs[1]), string(bs[2]))
	default:
		err = fmt.Errorf("unsupport")
	}
	return err
}
