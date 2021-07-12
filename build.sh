#!/bin/bash

##
# Builds binaries of rspace-cli for 3 amd64 architectures for win,max,linux
# Requires single argument, a version numbepr
# 2nd argument can be 'b' (build only, default); 'p' (publish only) or 'bp' (build and publish)
# 
# The Version number is set into cmd/version.go
# 
# 'build' rrequires a pre-created distribution folder  'dist' with 3 subfolders win.mac,linux
#
##

VERSION=$1
if [[ -z "VERSION" ]]; then
    echo "Version number required"
    exit 1
fi
DIST_DIR=dist
BUILD_LOG=$DIST_DIR/build.log


function build {
   ## set the version into version command.
   #sed -i -e  "s/\(var rsVersion = \)\S\+/\1\"$VERSION\"/" cmd/version.go
   # no uses go16 embed 
   echo "$version" > cmd/version.txt

    ### remove any old builds 
    rm -rf $DIST_DIR 
    echo "Building Windows amd64"
    env GOOS=windows GOARCH=amd64 go build -o $DIST_DIR/rscli-win-amd64-$VERSION.exe
    echo "Building Mac amd64"

    env GOOS=darwin GOARCH=amd64 go build -o $DIST_DIR/rscli-macosx-amd64-$VERSION
    echo "Building Linux amd64"
    env GOOS=linux GOARCH=amd64 go build -o $DIST_DIR/rscli-linux-amd64-$VERSION
    echo "Done, appending to build log $BUILD_LOG"
    echo "$VERSION built on $(date)" >> $BUILD_LOG
}

build
