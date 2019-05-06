#!/bin/bash

# Build the binary usimg docker (and move it to ./pi-gen/stage2/04-raspberry-converter/)
docker rm -v pigen_work
cp Dockerfile ..
cp goDeps.sh ..
docker build -t raspberryconvertercompiler ..
rm ../Dockerfile
rm ../goDeps.sh
docker run --rm -v $(pwd)/pi-gen/stage2/04-raspberry-converter:/out raspberryconvertercompiler

# Build the image
cd pi-gen
./build-docker.sh
mv deploy/*.zip ..
rm -rf deploy
