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

import "testing"

func Test_pathHandler(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				s: "/..//test/",
			},
			want: "/test/",
		},
		{
			name: "case2",
			args: args{
				s: "/../../test",
			},
			want: "/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathHandler(tt.args.s); got != tt.want {
				t.Errorf("pathHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
