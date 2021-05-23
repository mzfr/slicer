#!/bin/bash

env GOOS=windows GOARCH=amd64 go build
tar -cf slicer-windows-amd64.tar slicer.exe LICENSE readme.md config.yml

rm slicer.exe

env GOOS=darwin GOARCH=amd64 go build
tar -cf slicer-darwin-amd64.tar slicer LICENSE readme.md config.yml

rm slicer

go build
tar -cf slicer-linux.tar slicer LICENSE readme.md config.yml

rm slicer


