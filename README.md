# Fabric_based_application
基于Facric Hyperledger fabric-sdk-go区块链应用

# 环境
Ubuntu18.04及以上
GoLang(需要go版本1.13以上，最新版本兼容)
Docker
装docker-compose(开发时使用1.27.4)

# 安装使用
使用fabric-sdk-go前请先能够搭建起fabric网络，本链码基于fabric1.4及以上编写。

# 下载项目所需要的依赖
#进入项目目录下为项目文件夹赋权限：

#生成go.mod文件

go mod init   该命令会在该文件夹下生成一个go.mod和go.sum文件 正常输出类似于(go: creating new go.mod: module github.com/xxxxx)

#下载所需依赖（执行前将modelworker/chaincode中main.go和mwCC.go中import中的两行类似”github……”的两行注释掉，执行完成后再注释回来)

go mod vendor 该命令会在该文件夹下生成一个vendor文件夹，包含了所需依赖。

go env -w GOPROXY=https://goproxy.cn

#下载完成后即可进行测试：

#进入fixtures启动网络：(启动过就不需要再次启动)

cd fixtures

docker-compose -f hostn.yaml up -d

#编译并且运行：

#返回上一级有main.go的目录
#编译：

go build 会生成一个和当前文件夹同名的可执行文件，即modelworker

#运行

./modelworker

#运行成功提示：Fabric SDK初始化成功，通道已成功创建，peers 已成功加入通道.等等
