#!/bin/bash

# cd ./fixtures
# docker-compose -f host05.yaml down -v
# docker rm -f -v `docker ps -a --no-trunc | grep "mwcc" | cut -d ' ' -f 1` 2>/dev/null
# docker rmi `docker images --no-trunc | grep "mwcc" | cut -d ' ' -f 1` 2>/dev/null
# echo "环境清理成功"

cd ./fixtures
docker-compose -f host05.yaml up -d
echo "网络启动"
cd ../

# cd ../
# rm modelworker
# go build
# echo "二进制文件创建成功"

# ./modelworker
