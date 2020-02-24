#!/bin/bash

# Get dependencies
go get -u -v github.com/akamensky/argparse
go get -u -v github.com/jlaffaye/ftp
go get -u -v github.com/mattn/go-sqlite3
go get -u -v golang.org/x/crypto/ssh

# Build static binary
go build
