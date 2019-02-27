/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : vim.go
#   Created       : 2019-01-30 14:11:31
#   Describe      :
#
# ====================================================*/
package cmd

import (
	"os"
	"os/exec"
	"strings"
)

// Vim 使用vim修改etcd中的数据
func (r *Root) Vim(path string) error {

	path = strings.TrimSuffix(pathHandler(path), "/")

	tempDir := os.TempDir()

	r.touch(path)

	if err := r.download(path, tempDir); err != nil {
		return err
	}

	fn := fileName(path)
	localP := tempDir + fn
	defer func() {
		cmd := exec.Command("rm", localP)
		cmd.Run()
	}()

	cmd := exec.Command("vim", localP)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return r.upload(strings.TrimSuffix(path, fn), localP)
}
