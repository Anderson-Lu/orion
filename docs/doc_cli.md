**Orion-cli**是一个提效工具,通过**orion-cli**快速初始化服务项目

# 快速安装

```bash
go get github.com/Anderson-Lu/orion/tools/orion-cli
go install github.com/Anderson-Lu/orion/tools/orion-cli
```

这样就可以在`$GOPATH/bin`目录下看到`orion-cli`了

# 通过CLI创建项目模版(new)

> orion-cli new [项目目录路径] [项目模块名]

```bash
> orion-cli new demo github.com/cc

[orion-cli] Preparing ...
[orion-cli] Checking ...
[orion-cli] Create folder: demo
[orion-cli] create dir: demo/cmd
[orion-cli] create file: demo/cmd/main.go
[orion-cli] create dir: demo/proto
[orion-cli] create dir: demo/service
[orion-cli] create file: demo/service/service.go
[orion-cli] create dir: demo/config
[orion-cli] create file: demo/config/config.go
[orion-cli] create file: demo/go.mod
[orion-cli] Excute succ! please run `go mod tidy` to fix reference issues
```

目录结构如下:

```shell
.
├── cmd            #构建目录
│   └── main.go    #程序运行目录
├── config         #配置目录
│   └── config.go  #配置信息
├── go.mod         #go module初次不会生成go.sum,执行go mod tidy即可
├── Makefile       #构建脚本
├── proto          #pb协议目录,不是必须的,结合Makefile
└── service        #服务目录
    └── service.go #服务逻辑实现
```

更多细节可以通过`orion-cli --help`查看.