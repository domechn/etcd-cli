# ETCD-CLI

[中文文档](./README_zh.md)

## After connecting to the etcd service using the etcd-cli, you can use the common Linux commands to manipulate the data in the etcd

```bash
etcd-cli -s 127.0.0.1 -p 2380
```

### Supported commands

- cd
- ls
- mkdir
- touch
- rm
- mv
- cp
- pwd
- cat

`(additional parameters, such as -f, -r, -p, etc., are not supported)`

**you can even use `vim` to modify the etcd * * can be translated into text files**

It is important to note `rm`, `mv`, `cp` these commands in operation is needed in the add **"/"**, is used to distinguish the folder or file

```bash
# Delete the entire dir folder
rm dir/
# Delete the file
rm file
```

The same as `cp`,`mv`

![option](./images/option.gif)

### Additional commands supported

- upload [etcd-path] [local-path] Upload local **files** to the path specified by etc
- download [etcd-path] [local-path] Download the **file** specified in the etcd to local

If you are using upload or download in the connection state, the local-path needs to write the absolute path

You can also use etcd-cli directly, where local-path can use a relative path

```bash
etcd-cli -s 127.0.0.1 -p 2380 download /etcd-path/testfile ./
```

### How to install

First of all,install the go locale

```bash
go get github.com/hiruok/etcd-cli
cd $GOPATH/bin
./etcd-cli
```
