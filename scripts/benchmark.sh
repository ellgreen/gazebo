#!/usr/bin/env bash

set -e

package=$1

if [[ -z $package ]]
then
  package=./...
fi

go test -bench=. -run=^$ $package
