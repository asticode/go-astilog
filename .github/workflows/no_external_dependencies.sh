#!/bin/sh

if [[ $(go list -m all) =~ ^github.com/asticode/go-astilog[[:space:]]+github.com/asticode/go-astikit[[:space:]]+v[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+$ ]]; then
    echo "cheers"
else
    echo "This repo doesn't allow other dependencies than astikit"
    exit 1
fi