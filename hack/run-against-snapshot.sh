#!/usr/bin/env bash

SNAPSHOT_FILEPATH=${SNAPSHOT_FILEPATH:-$1}

if [[ ! -f $SNAPSHOT_FILEPATH ]]; then
    echo "Cannot find snapshot file. Please call $0 with path to snapshot or set SNAPSHOT_FILEPATH envvar."
    exit 1
fi

root_dir=$(cd "$(dirname "$0")/.."; pwd)
ghwc_image_name="ghwc"
local_git_version=$(git describe --tags --always --dirty)
IMAGE_VERSION=${IMAGE_VERSION:-$local_git_version}

snap_tmp_dir=$(mktemp -d -t ghw-snap-test-XXX)

echo "extracting snapshot $SNAPSHOT_FILEPATH to $snap_tmp_dir ..."
tar -xf $SNAPSHOT_FILEPATH -C $snap_tmp_dir

echo "building Docker image with ghwc ..."

docker build -f $root_dir/Dockerfile -t $ghwc_image_name:$IMAGE_VERSION $root_dir

echo "gathering symlink dirs in extracted snapshot to host volume mount ..."

linkdirs=""
for dirlink in `find -L $snap_tmp_dir -xtype l`; do
    sourcelink=$( readlink $dirlink )
    sourcelink=$( echo $sourcelink | sed "s/\/tmp\/ghw-snapshot[[:digit:]]\+/$tmp_snap_dir/")
    target=$( echo "$dirlink" | sed "s/\/tmp\/ghw-snap-test-.../\/host/")
    linkdirs+=" --mount type=bind,source=$sourcelink,destination=$target"
done

echo "running ghwc Docker image with volume mount to snapshot dir ..."

docker run -it -v $snap_tmp_dir:/host $linkdirs -e GHW_CHROOT="/host" $ghwc_image_name:$IMAGE_VERSION 
