#!/usr/bin/env bash

set -e

package=$1

if [[ -z $package ]]
then
  package=./...
fi

coverprofile="/tmp/go-cover.$$.tmp"

go test -v -coverpkg=$package -coverprofile=$coverprofile $package
go tool cover -html=$coverprofile
unlink $coverprofile
