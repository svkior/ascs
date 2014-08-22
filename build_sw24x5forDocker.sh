mkdir -p ./distout/docker

go build  -o ./distout/docker/main_cli main_cli.go
cp ./conf/lnx_switch24_docker.json ./distout/docker/
cd ./distout/docker
sudo docker build -t svkior/switcher24x7 .
CID=$(sudo docker run -d svkior/switcher24x7)

echo Press enter to quit
read

sudo docker stop $CID

