
BUILD_VERSION=`date -u +.%Y%m%d%.H%M%S`
MAJOR_VERSION="0.002"
DISTOUT=./distout


VERSION="$MAJOR_VERSION$BUILD_VERSION"
echo Building $VERSION
rm -rf $DISTOUT
mkdir -p $DISTOUT
echo building MAC64 
go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog_mac64 ffprog.go
echo building linux_386
GOARCH=386 GOOS=linux go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog_lnx32 ffprog.go
echo building linux_amd64
GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog_lnx64 ffprog.go
echo building windows32
GOARCH=386 GOOS=windows go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog.exe ffprog.go
echo building windows_amd64
GOARCH=amd64 GOOS=windows go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog64.exe ffprog.go

echo building linux_386
GOARCH=386 GOOS=linux go build -o ./distout/ffprog_gui_lnx32 ffprog_gui.go
