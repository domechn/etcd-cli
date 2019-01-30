# ETCD-CLI

## 使用etcd-cli连接etcd服务之后，可以使用常用的linux命令来操作etcd中的数据

```bash
etcd-cli -s 127.0.0.1 -p 2380
```

### 支持命令

- cd [path]
- ls [path]
- mkdir [path]
- touch [path]
- rm [path] [path] ...
- mv [source] [target]
- cp [source] [target]
- pwd
- cat [path]

**你甚至可以使用`vim`来修改etcd中可以被翻译成文本的文件**

需要注意的是`rm`、`mv`、`cp`这些命令在操作时需要在后加上 **"/"**,用来区分是文件夹还是文件

### 额外支持的命令

- upload [etcd-path] [local-path] 上传本地**文件**到etc指定的路径
- download [etcd-path] [local-path] 下载etcd中指定的**文件**到本地

如果在连接状态下使用upload或者download，local-path需要写绝对路径

也可以直接使用etcd-cli，这是local-path可以使用相对路径

```bash
etcd-cli -s 127.0.0.1 -p 2380 download /etcd-path/testfile ./
```

### 如何安装

安装go语言环境

```bash
go get github.com/hiruok/etcd-cli
cd $GOPATH/src/github.com/hiruok/etcd-cli/cmd
go build -o etcd-cli
mv etcd-cli $GOPATH/bin
```

### 请帮助完善ETCD-CLI