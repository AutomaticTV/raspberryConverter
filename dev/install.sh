#!/bin/bash

# Build the binary
sudo systemctl stop raspberryConverter.service
cd ..
packr build -v

# Create the service that will wun the binary on boot
# 1. Copy files to apropiate directories
sudo cp $GOPATH/src/raspberryConverter/build/pi-gen/stage2/04-raspberry-converter/raspberryConverter.service /etc/systemd/system/raspberryConverter.service
sudo cp $GOPATH/src/raspberryConverter/raspberryConverter /etc/systemd/raspberryConverter
# 2. Register the service
sudo systemctl enable raspberryConverter.service

# Disable login
sudo systemctl disable getty@tty1.service

# Silent boot
sudo sed -i 's/ quiet splash loglevel=0 logo.nologo vt.global_cursor_default=0//g' /boot/cmdline.txt
sudo sed -i '$ s/$/ quiet splash loglevel=0 logo.nologo vt.global_cursor_default=0/' /boot/cmdline.txt

# Full screen
sudo sed -i 's/#disable_overscan=1/disable_overscan=1/g' /boot/config.txt

# Disable rainbow splash screen
sudo sed -i 's/\ndisable_splash=1//g' /boot/config.txt
sudo sed -i '$ s/$/\ndisable_splash=1/' /boot/config.txt

# Custom splash screen
# 1. Copy files to apropiate directories
sudo cp $GOPATH/src/raspberryConverter/build/pi-gen/stage2/04-raspberry-converter/splashScreen.service /etc/systemd/system/splashScreen.service
sudo cp $GOPATH/src/raspberryConverter/build/pi-gen/stage2/04-raspberry-converter/splash.png /opt/splash.png
# 2. Register the service
sudo systemctl enable splashScreen.service

# Go back to dev
cd dev
