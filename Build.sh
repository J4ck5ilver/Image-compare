#!/bin/bash
set -e

TARGETS=("linux/amd64" "windows/amd64")

OUTPUT_DIR="./build"
mkdir -p $OUTPUT_DIR

for TARGET in "${TARGETS[@]}"; do
    OS=$(echo $TARGET | cut -d'/' -f1)
    ARCH=$(echo $TARGET | cut -d'/' -f2)
    OSARCH="${OS}_${ARCH}"
    EXTENSION=""
    if [ "$OS" == "windows" ]; then
        EXTENSION=".exe"
    fi
    echo "Building $OS/$ARCH..."

    cd compare/src
    echo "Building compare..."
    GOOS=$OS GOARCH=$ARCH go build -o ../../$OUTPUT_DIR/$OSARCH/Compare$EXTENSION
    cd ../..
    cd filter/src
    echo "Building filter..."
    GOOS=$OS GOARCH=$ARCH go build -o ../../$OUTPUT_DIR/$OSARCH/Filter$EXTENSION
    cd ../..
    cd browser/src
    echo "Building browser..."
    GOOS=$OS GOARCH=$ARCH go build -o ../../$OUTPUT_DIR/$OSARCH/Browser$EXTENSION
    cd ../..
done

cd $OUTPUT_DIR
for folder in *; do
    zip "${folder}.zip" "$folder"
done
cd -