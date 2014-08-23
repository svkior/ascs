tmpdir=$(mktemp -d)
pushd $tmpdir

BUSYBOX=`which busybox`

mkdir -p ./bin
cp $BUSYBOX ./bin/

# Create symbolic links back to busybox
for i in $(./bin/busybox --list);do
    ln -s /bin/busybox ./bin/$i
done

# Create container
tar -c . | docker import - svkior/barebusybox

# Go back to old pwd
popd

rm -rf $tmpdir