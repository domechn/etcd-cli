/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : etcd_test.go
#   Created       : 2019-01-29 11:49:23
#   Describe      :
#
# ====================================================*/
package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/domgoer/etcd-cli/pkg/store"
)

var st, _ = New([]string{"localhost:2379"}, nil)

func TestNew(t *testing.T) {
	_, err := New([]string{"localhost:2379"}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestEtcd_List(t *testing.T) {
	gf, err := st.List(context.Background(), "foo2/12/23")
	if err != nil {
		t.Error(err)
	}
	for _, v := range gf {
		fmt.Println(string(v.Value))
	}
}

func TestGetAndPut(t *testing.T) {
	err := st.Put(context.Background(), "fb", []byte("test"), &store.WriteOptions{TTL: time.Second})
	if err != nil {
		t.Error(err)
	}
	gf, err := st.Get(context.Background(), "fb")
	if err != nil {
		if err.Error() != "key not exsit" {
			t.Error(err)
		}
	} else {
		fmt.Println(string(gf.Value))
	}
}

func TestEtcd_Delete(t *testing.T) {
	err := st.Put(context.Background(), "fb-delete", []byte("test"), nil)
	if err != nil {
		t.Error(err)
	}
	gf, err := st.Get(context.Background(), "fb-delete")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(gf.LastIndex)
	err = st.Delete(context.Background(), "fb-delete")
	if err != nil {
		t.Error(err)
	}
	gf, err = st.Get(context.Background(), "fb-delete")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(gf.LastIndex)
}

func TestEtcd_AtomicPut(t *testing.T) {
	r, _ := st.Get(context.Background(), "foo2")
	r.LastIndex = 11
	// &{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:190 raft_term:5  true [response_put:<header:<revision:190 > > ]}
	st.AtomicPut(context.Background(), "foo2", []byte("5724"), r, nil)
}

func TestEtcd_AtomicDelete(t *testing.T) {
	// r.LastIndex = 11
	// &{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:190 raft_term:5  true [response_put:<header:<revision:190 > > ]}
	for i := 0; i < 100; i++ {
		go func() {
			r, _ := st.Get(context.Background(), "foo2")
			rs, _, _ := st.AtomicDelete(context.Background(), "foo2", r)
			if !rs {
				fmt.Println("触发cas")
			}

		}()
	}
}

func BenchmarkEtcd_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.Get(context.Background(), "foo2")
		}
	})
}

func BenchmarkEtcd_Put(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.Put(context.Background(), "foo2", []byte("32"), nil)
		}
	})
}

func BenchmarkEtcd_Delete(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.Delete(context.Background(), "foo2")
		}
	})
}

func BenchmarkEtcd_AtomicPut(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.AtomicPut(context.Background(), "foo2", []byte("232"), &store.KVPair{
				Key:       "foo2",
				Value:     []byte("235"),
				LastIndex: 123,
			}, nil)
		}
	})
}

func BenchmarkEtcd_AtomicDelete(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.AtomicDelete(context.Background(), "foo2", &store.KVPair{
				Key:       "foo2",
				Value:     []byte("235"),
				LastIndex: 123,
			})
		}
	})
}
