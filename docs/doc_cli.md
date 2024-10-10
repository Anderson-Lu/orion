**Orion-cli**是一个提效工具,通过**orion-cli**快速初始化服务项目

# 安装

```bash
go get github.com/Anderson-Lu/orion/tools/orion-cli
go build github.com/Anderson-Lu/orion/tools/orion-cli
```

这样就可以在`$GOPATH/bin`目录下看到`orion-cli`了

# 通过CLI创建项目模版(new)

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

