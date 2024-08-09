**English**|[简体中文](README.zh.md)

# govm

[![Go Reference](https://pkg.go.dev/badge/github.com/Open-Source-CQUT/govm.svg)](https://pkg.go.dev/github.com/Open-Source-CQUT/govm)
![Static Badge](https://img.shields.io/badge/go-1.22.5-blue)

govm is a command line tool for managing local go versions, which allows you to switch between different go versions more easily and quickly, and focus more on development work. It is a small tool that I wrote based on my usual usage habits and other similar open source tools. Since it is written in pure go, it can support mainstream windows, linux, and macos well.

## Installation

### Download

If you have a go environment and the version is greater than go1.16, you can use go install to install

```bash
$ go install github.com/Open-Source-CQUT/govm/cmd/govm@latest
```

Or you can download the latest binary file for the corresponding platform in [Release](https://github.com/Open-Source-CQUT/govm/releases). Currently, only releases for windows, macos, and linux are provided.

### linux

Install the govm file to the `/var/lib/govm` directory, and then link it to `/usr/local/bin`

```bash
$ ln -s /var/lib/govm/govm /usr/local/bin/govm
```

Check whether govm is available

```bash
$ govm version
govm versoin v1.0.0 linux/amd64
```

Use the install command to download the latest version

```bash
$ sudo govm install --use
```

Add the following content to `$HOME/.bashrc`

```bash
eval "$(govm profile -s --shell=bash)"
```

After re-login to the shell, test whether the go environment is available

```bash
$ go version
go version go1.22.5 linux/amd64
```

### windows

Add the location of `govm.exe` to the PATH system variable, and then confirm whether govm is available

```bash
$ govm version
govm versoin v1.0.0 windows/amd64
```

**gitbash**

Add the following content to the `%HOME/.bashrc` file

```bash
eval "$(govm profile -s --shell=gitbash)"
```

**powershell**

Add the following file to the `$env:USERPROFILE\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1` file. If the file does not exist, create it manually

```powershell
govm profile -s --shell=powershell | Out-String | Invoke-Expression
```

Re-login to the shell and test whether the go environment is available

```bash
$ go version
go version go1.22.5 windows/amd64
```

### macos

Install the govm binary file to the `/var/lib/govm` directory, and then link it to the `/usr/local/bin` directory

```bash
$ ln -s /var/lib/govm/govm /usr/local/bin/govm
```

Check whether govm is available

```bash
$ govm version
govm versoin v1.0.0 darwin/amd64
```

Use the install command to download the latest version

```bash
$ sudo govm install --use
```

Add the following content to the `$HOME/.zshrc` file

```bash
eval "$(govm profile -s --shell=bash)"
```

Re-login to the shell and test whether the go environment is available

```bash
$ go version
go version go1.22.5 darwin/amd64
```

### Other platforms

If you are a user of other platforms, go to [Go supported platforms](https://github.com/golang/go/blob/master/src/cmd/dist/build.go#L1727) to check whether your platform is supported, and then follow the steps below to compile.

First clone the source code to your local

```bash
$ git clone https://github.com/Open-Source-CQUT/govm.git
```

Switch to a specific version

```bash
$ git checkout tags/v1.0.0
```

Make sure you have the go compiler and make installed locally, and then execute your os and arch as parameters, as shown below

```bash
$ make build mode=release os=linux arch=amd64
```

After the compilation is completed, the compiled binary file will be generated in the `bin/release/` directory of the current project. Execute the following command to check whether it runs normally. The following output indicates that the compilation is successful.

```bash
$ ./govm version govm version untag linux/amd64
```

## Commands

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

govm has a total of 10 commands, most of which are very simple. The following is a simple demonstration of the use of the main commands.

### search

Search for available go versions. You can use regular expressions to match. By default, the first 20 items are sorted from high to low by version.

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

Search for a specific version 

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

Install the specified go version. If no parameters are specified, the latest version will be installed

```bash
$ govm install
Fetch go1.22.6 from https://dl.google.com/go/go1.22.6.windows-amd64.zip
Downloading go1.22.6.windows-amd64.zip 100% |█████████████████████████████████████| (76/76 MB, 34 MB/s) [2s]
Extract go1.22.6.windows-amd64.zip to local store
Remove archive from cache
Version go1.22.6 installed
```

Install and set to use version

```bash
$ sudo govm install 1.20.14
Fetch go1.20.14 from https://dl.google.com/go/go1.20.14.windows-amd64.zip
Downloading go1.20.14.windows-amd64.zip 100% |████████████████████████████████████| (114/114 MB, 32 MB/s) [3s]
Extract go1.20.14.windows-amd64.zip to local store
Remove archive from cache
Version go1.20.14 installed
Use go1.20.14 now
```

### use

Set an installed version as the used version

```bash
$ govm use 1.22.5
Use go1.22.5 now
```

### list

View the locally installed version

```bash
$ govm list
go1.22.6 (*)
go1.22.5
go1.22.3
go1.22.1
go1.21rc2
```

### uninstall

Uninstall a specific version

```bash
$ sudo govm uninstall 1.22.5
Version 1.22.5 uninstalled
```

### More

For more help information, please view it through `govm command help`

## Configuration

The configuration file of govm is stored in `$HOME/.govm/config.toml` in all systems. You can view the configuration by the following command

```bash
$ govm config
listapi=https://go.dev/dl/?mode=json&include=all
mirror=https://dl.google.com/go/
proxy=(system proxy)
install=/home/username/.govm/store/
```

### Mirror

The default download mirror of govm is to use the go official website. Chinese users are recommended to use the latter two

- Google: https://dl.google.com/go/, default
- Alibaba Cloud: https://mirrors.aliyun.com/golang/
- Nanjing University: https://mirrors.nju.edu.cn/golang/

Use the following command to modify the mirror.

```bash
$ govm cfg -w mirror=https://mirrors.aliyun.com/golang/
```

### Version list

The default version list uses the API provided by go officially

```
https://go.dev/dl/?mode=json&include=all
```

Modify it according to the following command

```bash
$ govm cfg -w listapi=your_cdn
```

### Proxy

The system proxy is used by default, and you can also specify the proxy manually, modify it with the following command

```bash
$ govm cfg -w proxy=your_proxy
```

### Installation path

The default storage location is in the `.govm/store/` directory, modify it with the following command

```bash
$ govm cfg -w install=new_pos
```