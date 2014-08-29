
BUILD_VERSION=`date -u +.%Y%m%d%.H%M%S`
MAJOR_VERSION="0.002"

VERSION="$MAJOR_VERSION$BUILD_VERSION"
echo Building $VERSION
mkdir -p ./distout
go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog_mac64 ffprog.go
GOARCH=386 GOOS=linux go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog_lnx32 ffprog.go
GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog_lnx64 ffprog.go
GOARCH=386 GOOS=windows go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog.exe ffprog.go
GOARCH=amd64 GOOS=windows go build -ldflags "-X main.version $VERSION" -o ./distout/ffprog64.exe ffprog.go
