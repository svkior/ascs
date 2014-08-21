WORK=./distout/docker/work
mkdir -p $WORK

go build  -o $WORK/main_cli main_cli.go
cp ./conf/lnx_switch24.json $WORK

