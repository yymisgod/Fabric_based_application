cd ./fixtures
docker-compose -f host05.yaml down -v
docker rm -f -v `docker ps -a --no-trunc | grep "mwcc" | cut -d ' ' -f 1` 2>/dev/null
docker rmi `docker images --no-trunc | grep "mwcc" | cut -d ' ' -f 1` 2>/dev/null
echo "环境清理成功"
cd ../
# docker ps -a
# docker images