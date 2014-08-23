mkdir -p ./distout/docker

go build  -o ./distout/docker/work/main_cli main_cli.go
cp ./conf/lnx_switch24_docker.json ./distout/docker/work/
pushd ./distout/docker
docker build -t svkior/switcher24x7 .
CID=$(docker run -d svkior/switcher24x7)

popd

echo docker logs -f $CID >current_log.sh
chmod +x current_log.sh

echo docker stop $CID >current_stop.sh
chmod +x current_stop.sh

echo Deployed with CID $CID
echo Press enter to quit
read

./current_stop.sh

rm current_log.sh
rm current_stop.sh

docker rm $CID

