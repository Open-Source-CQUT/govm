[English](README.en.md)|**简体中文**

# govm

govm是一个用于管理本地go版本的命令行工具，可以让你更简单和快速地切换不同的go版本，将更多注意力放在开发工作上。它是我结合平时使用习惯和借鉴了其他类似的开源工具而编写出来的一个小工具，由于它是纯go编写，所以应该能支持大部分的主流平台。



## 安装



### 下载

如果你拥有go环境，并且版本大于go1.16，可以采用go install来安装

```bash
$ go install github.com/Open-Source-CQUT/govm/cmd/govm@latest
```

或者可以在[Release](https://github.com/Open-Source-CQUT/govm/releases)中下载对应平台的最新版二进制文件，目前仅提供windows，macos，linux三个平台的发行版。



### 编译

如果你是其他平台的用户，前往[Go supported platforms](https://github.com/golang/go/blob/master/src/cmd/dist/build.go#L1727)查阅具体支持哪些平台，然后按照下面的步骤自行编译。

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

然后会在当前项目的`./bin/release/`目录下生成编译好的二进制文件，执行如下命令查看是否正常运行，出现如下输出表示编译成功。

```bash
$ ./govm -v
govm version v1.0.0 linux/amd64
```



## 使用



### 快速开始

使用install命令下载最新版

```bash
$ govm install
```

在linux，macos系统下设置环境变量，此方法适用于gitbash for windows

```bash
$ govm profile >> $HOME/.profile && source $HOME/.profile
```

在powershell中设置环境变量，需新开powershell生效

```powershell
PS C:\Users\Administrator> setx /M PATH "$env:PATH;$env:USERPROFILE\.govm\root\go\bin"
```

测试go是否可用

```bash
$ go version
go version go1.22.5 linux/amd64
```



### 镜像

govm的默认下载镜像是使用go官网，中国用户建议使用后两个

- 谷歌：`https://dl.google.com/go/`

- 阿里云：`https://mirrors.aliyun.com/golang/`

- 南京大学：`https://mirrors.nju.edu.cn/golang/`

中科大虽然也有go镜像，但是会报403，不建议使用。



### 版本列表

默认的版本列表使用的是go官方提供的API

```
https://go.dev/dl/?mode=json&include=all
```

国内用户应该是比较难访问的，不过这是一个可配置项，可以自行修改。

