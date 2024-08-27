# 依赖

```shell
go get go.bug.st/serial
```

# 测试

```shell
# gopkg 项目目录
## 测试所有包
make 

## 只测试串口模块
go test -v -count=1 ./pkg/uart
```
