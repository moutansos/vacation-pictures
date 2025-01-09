#! /bin/bash

go env -w CGO_ENABLED=1
go env -w CC=musl-gcc

export CGO_LDFLAGS="-static -Wl,-unresolved-symbols=ignore-all"

# pulumi refresh
pulumi up
