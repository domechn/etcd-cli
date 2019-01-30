/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : watch.go
#   Created       : 2019-01-29 16:41:41
#   Describe      :
#
# ====================================================*/
package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/hiruok/etcd-cli/pkg/store/etcd"

	"github.com/hiruok/etcd-cli/pkg/store"
)

var (
	ss, _ = etcdstore.New([]string{"127.0.0.1:2381"}, nil)
)

func TestClient_WatchConfByEnv(t *testing.T) {
	type fields struct {
		s        store.Store
		fileType string
	}
	type args struct {
		ctx    context.Context
		key    string
		stopCh <-chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    <-chan []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "watch_test_case_1",
			fields: fields{
				s:        ss,
				fileType: "yaml",
			},
			args: args{
				ctx:    context.Background(),
				key:    "/test-watch",
				stopCh: make(chan struct{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				s:        tt.fields.s,
				fileType: "yaml",
			}
			sss, err := c.WatchConfByEnv(tt.args.ctx, tt.args.key, tt.args.stopCh)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.WatchConfByEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			go func() {
				for {
					a := <-sss
					fmt.Println(string(a))
				}
			}()
		})
	}
}
