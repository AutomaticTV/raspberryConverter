#!/bin/bash

# Clone pi-gen and copy
git clone https://github.com/RPi-Distro/pi-gen.git
mkdir pi-gen/stage2/04-raspberry-converter
cp -r 04-raspberry-converter pi-gen/stage2/
cp config pi-gen

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
touch stage3/SKIP
touch stage4/SKIP
touch stage5/SKIP
rm stage4/EXPORT_IMAGE
rm stage4/EXPORT_NOOBS
rm stage5/EXPORT_IMAGE
rm stage5/EXPORT_NOOBS
docker-compose up -d
./build-docker.sh
mv deploy/*.zip ..

# Clean
cd ..
rm -rf pi-gen
