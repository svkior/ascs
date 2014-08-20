mkdir -p ./distout/docker

go build  -o ./distout/docker/main_cli main_cli.go
cp ./conf/lnx_switch24.json ./distout/docker/


