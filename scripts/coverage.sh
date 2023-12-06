#!/bin/bash

for d in $(go list ./...) 
do
go test -cover "$d" | grep -v ?
done