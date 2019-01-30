/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : path.go
#   Created       : 2019-01-29 19:27:41
#   Describe      :
#
# ====================================================*/
package utils

import "strings"

// Normalize the key for each store to the form:
//
//     /path/to/key
//
func Normalize(key string) string {
	return "/" + join(SplitKey(key))
}

// SplitKey splits the key to extract path informations
func SplitKey(key string) (res []string) {
	var path []string
	if strings.Contains(key, "/") {
		path = strings.Split(key, "/")
	} else {
		path = []string{key}
	}
	for _, p := range path {
		if p != "" {
			res = append(res, p)
		}
	}
	return
}

// join the path parts with '/'
func join(parts []string) string {
	return strings.Join(parts, "/")
}
