#!/bin/bash

GO111MODULE=on  \
CGO_CFLAGS="-Wl,--no-as-needed -w -D_XOPEN_SOURCE=600 -std=gnu99 -Werror=implicit-function-declaration -lpthread" \
go run -v test.go
