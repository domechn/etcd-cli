/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : helper.go
#   Created       : 2019-01-30 10:09:12
#   Describe      :
#
# ====================================================*/
package opts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hiruok/etcd-cli/pkg/stack"

	"gopkg.in/yaml.v2"

	"github.com/json-iterator/go"
)

var (
	// ErrInvalidParamNum 参数数量错误
	ErrInvalidParamNum = errors.New("invalidate param nums")

	// PWD 当前位置
	PWD = "/"

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func readFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, err
}

func validate(b []byte, path string) error {
	types := strings.Split(path, ".")
	typeN := types[len(types)-1]
	var mmap = make(map[string]interface{})
	switch typeN {
	case "yml", "yaml":
		return yaml.Unmarshal(b, &mmap)
	case "json":
		return json.Unmarshal(b, &mmap)
	}
	return nil
}

func fileName(s string) string {
	ss := strings.Split(s, "/")
	res := ss[len(ss)-1]
	if res == "" {
		return "unknowfile"
	}
	return res
}

func changeLine() {
	fmt.Println()
}

func pathHandler(s string) string {
	if !strings.HasPrefix(s, "/") {
		s = PWD + s
	}
	s = strings.TrimRight(s, "/")
	sk := stack.New()
	ss := strings.Split(s, "/")
	for _, v := range ss {
		if v == "" {
			sk.Flush()
			continue
		}
		if v == "." {
			continue
		}
		if v == ".." {
			sk.Pop()
			continue
		}
		sk.Push(&stack.Item{
			Value: []byte(v),
		})
	}
	bRes := sk.Value()
	var resList []string
	for _, v := range bRes {
		resList = append(resList, string(v.Value))
	}
	path := strings.Join(resList, "/") + "/"
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}
