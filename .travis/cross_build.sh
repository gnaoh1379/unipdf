#!/usr/bin/env bash

# Functions.
function info() {
    echo -e "\033[00;34mi\033[0m $1"
}

function fail() {
    echo -e "\033[00;31m!\033[0m $1"
    exit 1
}

function build() {
    goos=$1
    goarch=$2

    info "Building for $goos $goarch..."
    GOOS=$goos GOARCH=$goarch go build -o $goos_$goarch main.go
    if [[ $? -ne 0 ]]; then
        fail "Could not build for $goos $goarch. Aborting."
    fi
}

# Create build directory.
mkdir -p bin
cd bin

# Create go.mod
cat <<EOF > go.mod
module cross_build
require github.com/gnaoh1379/unipdf v3.0.0
EOF

echo "replace github.com/gnaoh1379/unipdf => $TRAVIS_BUILD_DIR" >> go.mod

# Create Go file.
cat <<EOF > main.go
package main

import (
	_ "github.com/gnaoh1379/unipdf/annotator"
	_ "github.com/gnaoh1379/unipdf/common"
	_ "github.com/gnaoh1379/unipdf/common/license"
	_ "github.com/gnaoh1379/unipdf/contentstream"
	_ "github.com/gnaoh1379/unipdf/contentstream/draw"
	_ "github.com/gnaoh1379/unipdf/core"
	_ "github.com/gnaoh1379/unipdf/core/security"
	_ "github.com/gnaoh1379/unipdf/core/security/crypt"
	_ "github.com/gnaoh1379/unipdf/creator"
	_ "github.com/gnaoh1379/unipdf/extractor"
	_ "github.com/gnaoh1379/unipdf/fdf"
	_ "github.com/gnaoh1379/unipdf/fjson"
	_ "github.com/gnaoh1379/unipdf/model"
	_ "github.com/gnaoh1379/unipdf/model/optimize"
	_ "github.com/gnaoh1379/unipdf/model/sighandler"
	_ "github.com/gnaoh1379/unipdf/ps"
	_ "github.com/gnaoh1379/unipdf/render"
)

func main() {}
EOF

# Build file.
for os in "linux" "darwin" "windows"; do
    for arch in "386" "amd64"; do
        build $os $arch
    done
done
