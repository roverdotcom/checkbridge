#!/bin/bash

set -eu

BIN_NAME="checkbridge"

DIST_DIR="$(pwd)/dist/out"
mkdir -p "${DIST_DIR}"

VERSION="$(git tag -l --points-at HEAD)"
if [[ -z "${VERSION}" ]] ; then
  echo "Release can only be run on a tagged commit."
  exit 2
fi

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "${BRANCH}" != "master" ]] ; then
  echo "Release can only be run from the master branch"
  exit 3
fi

if ! grep "$VERSION" README.md > /dev/null ; then
  echo "Readme not pointing downloads to $VERSION"
  exit 4
fi

for os in linux darwin; do
  for arch in amd64 arm64; do
    name="${BIN_NAME}-$VERSION.$os-$arch"
    path="${DIST_DIR}/${name}"
    echo "Building $name - $VERSION"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags "-X github.com/roverdotcom/checkbridge/cmd.Version=$VERSION" \
        -o $path
    gzip < $path > $path.gz
    cp "$path" "${DIST_DIR}/${BIN_NAME}"
    (cd ${DIST_DIR} && tar zcf $path.tar.gz "${BIN_NAME}" && rm "${BIN_NAME}")
    mkdir -p "$DIST_DIR/$os/$arch"
    cp "$path" "$DIST_DIR/$os/$arch/${BIN_NAME}"
  done
done
