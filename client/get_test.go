/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : get.go
#   Created       : 2019-01-29 16:26:41
#   Describe      :
#
# ====================================================*/
package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/hiruok/etcd-cli/pkg/store"
)

func TestClient_GetConfByEnv(t *testing.T) {
	type fields struct {
		s store.Store
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get_test_case_1",
			fields: fields{
				s: ss,
			},
			args: args{
				ctx: context.Background(),
				key: "/test-watch",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				s: tt.fields.s,
			}
			got, err := c.GetConfByEnv(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetConfByEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(string(got))
		})
	}
}
