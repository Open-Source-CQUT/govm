[English](README.md)|**简体中文**

# govm

[![Go Reference](https://pkg.go.dev/badge/github.com/Open-Source-CQUT/govm.svg)](https://pkg.go.dev/github.com/Open-Source-CQUT/govm)
![Static Badge](https://img.shields.io/badge/go-1.22.5-blue)

govm是一个用于管理本地go版本的命令行工具，可以让你更简单和快速地切换不同的go版本，将更多注意力放在开发工作上。它是我结合平时使用习惯和借鉴了其他同类的开源工具而编写出来的一个小工具，由于它是纯go编写，对于主流的windows，linux，macos都能良好支持。



## 安装



### 下载

如果你拥有go环境，并且版本大于go1.16，可以采用go install来安装

```bash
$ go install github.com/Open-Source-CQUT/govm/cmd/govm@latest
```

或者可以在[Release](https://github.com/Open-Source-CQUT/govm/releases)中下载对应平台的最新版二进制文件，目前仅提供windows，macos，linux三个平台的发行版。



### linux

将govm文件安装到`/var/lib/govm`目录下，再链接至`/usr/local/bin`

```bash
$ ln -s /var/lib/govm/govm /usr/local/bin/govm
```

查看govm是否可用

```bash
$ govm version
govm versoin v1.0.0 linux/amd64
```

使用install命令下载最新版

```bash
$ sudo govm install --use
```

将下列内容添加到`$HOME/.bashrc`中

```bash
eval "$(govm profile -s --shell=bash)"
```

重新登陆shell后测试go环境是否可用

```bash
$ go version
go version go1.22.5 linux/amd64
```



### windows

将`govm.exe`的位置添加到PATH系统变量中，然后确认govm是否可用

```bash
$ govm version
govm versoin v1.0.0 windows/amd64
```

**gitbash**

将下列内容添加到`%HOME/.bashrc`文件中

```bash
eval "$(govm profile -s --shell=gitbash)"
```

**powershell**

将下列文件添加到`$env:USERPROFILE\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1`文件中，如果不存在该文件就手动创建

```powershell
govm profile -s --shell=powershell | Out-String | Invoke-Expression
```

重新登陆shell后测试go环境是否可用

```bash
$ go version
go version go1.22.5 windows/amd64
```



### macos

将govm二进制文件安装到`/var/lib/govm`目录下，再链接至`/usr/local/bin`目录下

```bash
$ ln -s /var/lib/govm/govm /usr/local/bin/govm
```

查看govm是否可用

```bash
$ govm version
govm versoin v1.0.0 darwin/amd64
```

使用install命令下载最新版

```bash
$ sudo govm install --use
```

将下列内容添加到`$HOME/.zshrc`文件中

```bash
eval "$(govm profile -s --shell=bash)"
```

重新登陆shell后测试go环境是否可用

```bash
$ go version
go version go1.22.5 darwin/amd64
```



### 其他平台

如果你是其他平台的用户，前往[Go supported platforms](https://github.com/golang/go/blob/master/src/cmd/dist/build.go#L1727)查阅是否支持你的平台，然后按照下面的步骤编译。

首先将源代码克隆到本地

```bash
$ git clone https://github.com/Open-Source-CQUT/govm.git
```

切换到特定版本

```bash
$ git checkout tags/v1.0.0
```

确保你本地安装了go编译器和make，然后并将你的os和arch作为参数执行，示例如下

```bash
$ make build mode=release os=linux arch=amd64
```

编译完成后会在当前项目的`bin/release/`目录下生成编译好的二进制文件，执行如下命令查看是否正常运行，出现如下输出表示编译成功。

```bash
$ ./govm version
govm version untag linux/amd64
```



## 命令

```bash
$ govm -h
govm is a tool to manage local Go versions

Usage:
  govm [command]

Available Commands:
  clean       Clean local cache and redundant versions
  completion  Generate the autocompletion script for the specified shell
  config      Manage govm configs
  current     Show current using Go version
  help        Help about any command
  install     Install specified Go version
  list        List local installed Go versions
  profile     Show profile env
  search      Search available go versions from remote
  uninstall   Uninstall specified Go version
  use         Use specified Go version
  version     Show govm version

Flags:
  -h, --help      help for govm
  -s, --silence   Do not show any tip, warn, error

Use "govm [command] --help" for more information about a command.
```

govm总共由10个命令，大部分都很简单，下面简单演示下主要命令的使用。



### search

搜索可用的go版本，可以用正则进行匹配，默认从高到低按照版本进行排序显示前20条。

```bash
$ govm search
go1.22.6  	   69 MB
go1.22.5  	   69 MB
go1.22.4  	   69 MB
go1.22.3  	   69 MB
go1.22.2  	   69 MB
go1.22.1  	   69 MB
go1.22.0  	   69 MB
go1.21.13 	   67 MB
go1.21.12 	   67 MB
go1.21.11 	   67 MB
......
go1.21.1  	   67 MB
```

搜索特定的版本

```bash
$ govm search 1.18 -n 10
go1.18.10 	  142 MB
go1.18.9  	  142 MB
go1.18.8  	  142 MB
go1.18.7  	  142 MB
go1.18.6  	  142 MB
go1.18.5  	  142 MB
go1.18.4  	  142 MB
go1.18.3  	  142 MB
go1.18.2  	  142 MB
go1.18.1  	  142 MB
```



### install

安装指定的go版本，不指定任何参数时则安装最新版本

```bash
$ govm install
Fetch go1.22.6 from https://dl.google.com/go/go1.22.6.windows-amd64.zip
Downloading go1.22.6.windows-amd64.zip 100% |████████████████████████████████████████| (76/76 MB, 34 MB/s) [2s]
Extract go1.22.6.windows-amd64.zip to local store
Remove archive from cache
Version go1.22.6 installed
```

安装并设置为使用版本

```bash
$ sudo govm install 1.20.14
Fetch go1.20.14 from https://dl.google.com/go/go1.20.14.windows-amd64.zip
Downloading go1.20.14.windows-amd64.zip 100% |████████████████████████████████████████| (114/114 MB, 32 MB/s) [3s]
Extract go1.20.14.windows-amd64.zip to local store
Remove archive from cache
Version go1.20.14 installed
Use go1.20.14 now
```



### use

将某一个已安装的版本设置为使用版本

```bash
$ govm use 1.22.5
Use go1.22.5 now
```



### list

查看本地已安装的版本

```bash
$ govm list
go1.22.6 (*)
go1.22.5
go1.22.3
go1.22.1
go1.21rc2
```



### uninstall

卸载某一个特定版本

```bash
$ sudo govm uninstall 1.22.5
Version 1.22.5 uninstalled
```



### 更多

更多帮助信息请通过`govm command help`来查看



## 配置

govm的配置文件在所有系统中都存放在`$HOME/.govm/config.toml`中，通过如下命令可以查看配置

```bash
$ govm config
listapi=https://go.dev/dl/?mode=json&include=all
mirror=https://dl.google.com/go/
proxy=(system proxy)
install=/home/username/.govm/store/
```



### 镜像

govm的默认下载镜像是使用go官网，中国用户建议使用后两个

- 谷歌：https://dl.google.com/go/，默认
- 阿里云：https://mirrors.aliyun.com/golang/
- 南京大学：https://mirrors.nju.edu.cn/golang/

**中科大虽然也有go镜像，但是会报403，不推荐使用**。使用如下命令修改镜像

```bash
$ govm cfg -w mirror=https://mirrors.aliyun.com/golang/
```



### 版本列表

默认的版本列表使用的是go官方提供的API

```
https://go.dev/dl/?mode=json&include=all
```

中国用户应该是比较难访问的，不过这是一个可配置项，按照如下命令修改

```bash
$ govm cfg -w listapi=your_cdn
```



### 代理

默认情况下使用系统代理 ，也可以手动指定代理，使用如下命令修改

```bash
$ govm cfg -w proxy=your_proxy
```



### 安装路径

默认存放位置位于`.govm/store/`目录下，使用如下命令修改

```bash
$ govm cfg -w install=new_pos
```
