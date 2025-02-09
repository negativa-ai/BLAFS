#!/bin/bash

mkdir -p build
cd build
cmake ../..
make

cd ..
cp build/file_reader .

docker build -t file_reader .

rm -rf build
rm file_reader