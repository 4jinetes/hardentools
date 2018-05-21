#!/bin/bash

GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1 go get -u --ldflags '-s -w -extldflags "-static" -H windowsgui' github.com/lxn/win
GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1 go get -u --ldflags '-s -w -extldflags "-static" -H windowsgui' github.com/lxn/walk
go get github.com/akavel/rsrc
go get golang.org/x/sys/windows/registry
go get gopkg.in/Knetic/govaluate.v3

# check go version
GO_VERSION="$(go version)"
GO_VERSION="$(echo $GO_VERSION | awk '{print $3}')"
if [[ $GO_VERSION == "go1.10"* ]] || [[ $GO_VERSION == "go1.9"* ]] || [[ $GO_VERSION == "go1.8"* ]]; then
	$GOPATH/bin/rsrc -manifest harden.manifest -ico harden.ico -o rsrc.syso
	GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1 go build --ldflags '-s -w -extldflags "-static" -H windowsgui' -o hardentools.exe
else
	echo "Error: Build currently only works with go 1.8.x or go1.9.X"
fi
