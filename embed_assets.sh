#!/usr/bin/env bash

# Must have go-bindata installed (go get -u github.com/jteeuwen/go-bindata/...)

go-bindata -o ./sheet_file/bindata.go -pkg sheet_file ./base_sheet.xlsx
