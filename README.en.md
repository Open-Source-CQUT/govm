**English**|[简体中文](README.md)

# govm

GOVM is a command line tool for managing the local GO version, which allows you to switch different GO versions simpler and quickly and put more attention on development. It is a small tool written by my usual use habits and borrowing other similar open source tools. Since it is written in pure Go, it should be able to support most mainstream platforms.



## Install



### download

If you have a GO environment and the version is greater than Go1.16, you can use Go Install to install

`` `Bash
$ Go Install github.com/open-source-cqut/govm/govm@latest
`` `

Or you can download the latest version of the binary files of the corresponding platform in [Release] (https://github.com/open-source-cqut/Govm/releases). At present, only the distribution of Windows, Macos, and Linux.



### Compilation

If you are a user of other platforms, go to [Go SUPPORTED PLATFORMS] (https://github.com/golang/go/blob/master/cmd/dist/build.go1727), and then what platforms are the specific support support. Declarily compile according to the following steps.

First clone the source code to the local area

`` `Bash
$ git clone https://github.com/open-source-cqut/govm.git
`` `

Switch to a specific version

`` `Bash
$ git checkout tags/v1.0.0
`` `

Make sure you have installed the Go compiler and make locally, and then use your OS and Arch as a parameter. The example is as follows

`` `Bash
$ Make Build Mode = Release OS = Linux Arch = AMD64
`` `

Then generate a binary file that is compiled in the current project's `./Bin/release/` directory. The following commands are executed to see if it runs normally.

`` `Bash
$ ./govm -V
Govm Version v1.0.0 Linux/AMD64
`` `



## use



### Fast start

Use the Install command to download the latest version

`` `Bash
$ GOVM Install
`` `

Set up environmental variables under Linux, MacOS system, this method is suitable for gitbash for Windows

`` `Bash
$ GOVM Profile >> $ HOME/.profile && Source $ Home/.profile
`` `

Set up environmental variables in Powershell, you need to open PowerShell to take effect

`` `PowerShell
PS C: \ Users \ Administrator> Setx /M Path "$ ENV: Path; $ ENV: Userprofile \ .govm \ ROOT \ GO \ BIN" ""
`` `

Test whether it is available

`` `Bash
$ Go Version
Go Version Go1.22.5 Linux/AMD64
`` `



### mirror

GOVM's default download mirror is used to use Go official website. Chinese users suggest that the post -use two latter two

-Google: `https: // dl.google.com/Go/`

-Alibaba Cloud: `https: // mirrs.aliyun.com/Golang/`

-Nanjing University: `https: // mirrs.nju.edu.cn/Golang/`

Although the University of Science and Technology also has GO mirror, it will be reported 403, and it is not recommended.



### version list

The default version list uses the official API provided by Go

`` `
https://go.dev/dl/?mode=json&include=all
`` `

Chinese users should be more difficult to access, but this is a configurable item that can be modified by itself.