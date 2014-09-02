
BUILD_VERSION=`date -u +.%Y%m%d%.%M:%S`
DISTOUT=./distout

perl -i -pe 's/echo \d+\.\d+\.\K(\d+)/ $1+1 /e' _version

chmod +x _version
MAJOR_VERSION=`./_version`

OUT_FILE_SUFFIX="_$MAJOR_VERSION"
VERSION="$MAJOR_VERSION at $BUILD_VERSION"
echo Building $VERSION
echo FileName: $OUT_FILE_SUFFIX

rm -rf $DISTOUT
mkdir -p $DISTOUT
echo building MAC64 
go build -ldflags "-X main.version \"$VERSION\"" -o ./distout/ffprog_mac64$OUT_FILE_SUFFIX ffprog.go

echo building linux_386
GOARCH=386 GOOS=linux go build -ldflags "-X main.version \"$VERSION\"" -o ./distout/ffprog_lnx32$OUT_FILE_SUFFIX ffprog.go

echo building linux_amd64
GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version \"$VERSION\"" -o ./distout/ffprog_lnx64$OUT_FILE_SUFFIX ffprog.go

echo building windows32
GOARCH=386 GOOS=windows go build -ldflags "-X main.version \"$VERSION\"" -o ./distout/ffprog$OUT_FILE_SUFFIX.exe ffprog.go

echo building windows_amd64
GOARCH=amd64 GOOS=windows go build -ldflags "-X main.version \"$VERSION\"" -o ./distout/ffprog64$OUT_FILE_SUFFIX.exe ffprog.go
