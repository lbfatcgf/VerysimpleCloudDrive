#!/bin/bash
set -e -u

if [ -d "build" ]; then
    printf "create build\n"
else
    mkdir build
    printf "create build\n"
fi

if [ -d "build/vscd" ]; then
    printf "create build/vscd\n"
else
    mkdir build/vscd
    printf "create build/vscd\n"
fi

if [ -d "build/view" ]; then
    printf "create build/view\n"
else
    mkdir build/view
    printf "create build/view\n"
fi

if [ -d "build/view/index.html" ]; then
    printf "create build/view/index.html\n"
else
    cp view/index.html build/view/index.html
    printf "create build/view/index.html\n"
fi
go build -o build/VSCD   VerysimpleCloudDrive 