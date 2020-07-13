# STORE

## etcd

### 示例

```go

package main

import (
	"github.com/domgoer/gateway/pkg/store/etcd"
	"github.com/domgoer/gateway/pkg/store"
	"fmt"
	"time"
	"context"
	)

func main() {
    client, err := etcdStore.New([]string{"localhost:2379"}, nil)
    if err != nil {
    	    panic(err)
    }
    //get 获取单个键的值
    getRes , err := client.Get(context.Background(),"foo")
    fmt.Printf("%#v\n",getRes)
    //list 列出目录下所有键值
    listRes , err := client.List(context.Background(),"/usr/local")
    fmt.Printf("%#v\n",listRes)
    //put 添加或修改键
    err = client.Put(context.Background(),"foo",[]byte("testvalue"),&store.WriteOptions{
    	TTL:time.Second*20,
    })
    //delete 删除单个键
    err = client.Delete(context.Background(),"foo")
    //deletetree 删除目录下所有键
    err = client.DeleteTree(context.Background(),"/usr/local/foo")
    //watch and watchtree  观察某个键或者目录
    //观察目录时，当目录下一个值变化，会将整个目录返回
    stopCh := make(chan struct{})
    	watchRes, err := client.Watch(context.Background(),"foo2/23", stopCh)
    	go func() {
    		for {
    			s, ok := <-watchRes
    			if ok {
    				fmt.Println(string(s.KV.Value))
    			}
    		}
    	}()
    	time.Sleep(time.Second * 5)
    	stopCh <- struct{}{}
    	//AtomicPut and AtomicDelete
    	flag , np ,err := client.AtomicPut(context.Background(),"foo",[]byte("new-value"),getRes,nil)
    	if !flag{ 
    		fmt.Println(np)
    		//如果flag为false说明更新失败
    		//如果err为nil,就会返回foo最新的版本
    	}
    	
    	//lock and unlock
    	loc , err := client.NewLock(etcdStore.DefaultLockKey,nil)
    	loc.Lock()
    	loc.Unlock()
}
```

### Operation参数
```go
    package store
    
    import (
    	"crypto/tls"
    "time"
    )
    
    // WriteOptions put 和 atomicput 时使用
    type WriteOptions struct {
    	IsDir bool          //该键是否是目录，如果为true，会默认在key后面添加"/"
    	TTL   time.Duration //该键过期时间
    }
    
    // LockOptions newlock 时使用
    type LockOptions struct {
    	TTL time.Duration   // lock的过期时间，如果未填默认20s
    }
    
    // KVPair get时的返回值
    type KVPair struct {
    	Key       string    //键
    	Value     []byte    //值
    	LastIndex uint64    //最新的版本号
    }
    
    
   
    //https 配置
    type Config struct {
    	TLS               *tls.Config       //tls配置信息
    	ConnectionTimeout time.Duration     //连接超时时间
    	Username          string            //用户名密码必须同时存在，否则不生效
    	Password          string
    }

```

### 说明

测试环境 : macbook pro 2015 13寸 中配

```
 go test -bench=. -benchtime="3s"
 goos: darwin
 goarch: amd64
 pkg: github.com/domgoer/github.com/domgoer/gateway/pkg/store
 BenchmarkEtcd_Get-4                50000             88290 ns/op
 BenchmarkEtcd_Put-4                 2000           2726530 ns/op
 BenchmarkEtcd_Delete-4              2000           2616983 ns/op
 BenchmarkEtcd_AtomicPut-4           2000           2658369 ns/op
 BenchmarkEtcd_AtomicDelete-4        2000           2595424 ns/op
```

```
go test -bench=. -benchtime="10s"
goos: darwin
goarch: amd64
pkg: github.com/domgoer/gateway/pkg/store
BenchmarkEtcd_Get-4               200000             90198 ns/op
BenchmarkEtcd_Put-4                 5000           2714720 ns/op
BenchmarkEtcd_Delete-4              5000           2611590 ns/op
BenchmarkEtcd_AtomicPut-4           5000           2593163 ns/op
BenchmarkEtcd_AtomicDelete-4        5000           2627109 ns/op
```

 > 线程不安全