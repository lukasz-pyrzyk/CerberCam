#!/bin/sh
cd ..
protoc --go_out=. *.proto
