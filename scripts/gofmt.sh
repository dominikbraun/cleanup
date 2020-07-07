#!/bin/sh
set -e
set -x

! (go fmt ./... | read)