#!/bin/bash

for d in $(go list ./...)
do
go test -v "$d" | grep -v ?
done