#!/bin/bash

set -eu

BIN_NAME="checkbridge"

DIST_DIR="$(pwd)/dist/out"
mkdir -p "${DIST_DIR}"

VERSION="${1:-}"
if [[ -z "${VERSION}" ]] ; then
  echo "Need version specifier"
  exit 2
fi

ARCH="$(go env GOARCH)"

for os in linux darwin; do
    name="${BIN_NAME}-$VERSION.$os-$ARCH"
    path="${DIST_DIR}/${name}"
    echo "Building $name - $VERSION"
    GOOS=$os GOARCH=$ARCH CGO_ENABLED=0 go build \
        -o $path
    gzip < $path > $path.gz
    cp "$path" "${DIST_DIR}/${BIN_NAME}"
    (cd ${DIST_DIR} && tar zcf $path.tar.gz "${BIN_NAME}" && rm "${BIN_NAME}")
    mkdir -p "$DIST_DIR/$os"
    cp "$path" "$DIST_DIR/$os/${BIN_NAME}"
done
