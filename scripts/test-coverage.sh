#!/usr/bin/env bash

if [[ -z $1 ]]
then
  printf "Usage: %s PACKAGE\n" $0
  exit 1
fi

coverprofile="/tmp/go-cover.$$.tmp"

go test -coverpkg=$1 -coverprofile=$coverprofile $1 && go tool cover -html=$coverprofile && unlink $coverprofile
