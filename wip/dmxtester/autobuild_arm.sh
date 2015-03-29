
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

echo building linux_arm
GOARCH=arm GOARM=5 GOOS=linux go build -o ./distout/dmxtester$OUT_FILE_SUFFIX dmxtester.go
#GOARCH=arm GOARM=5 GOOS=linux go build -ldflags "-X main.version \"$VERSION\"" -o ./distout/dmxtester$OUT_FILE_SUFFIX dmxtester.go

ffprog -json ./build_5.json