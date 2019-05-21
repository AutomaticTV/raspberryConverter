#!/bin/bash -e

# COPY FILES TO THE OS FILESYSTEM
cp raspberryConverter.service "${ROOTFS_DIR}/etc/systemd/system/raspberryConverter.service"
cp raspberryConverter "${ROOTFS_DIR}/etc/systemd/raspberryConverter"
cp splashScreen.service "${ROOTFS_DIR}/etc/systemd/system/splashScreen.service"
cp splash.png "${ROOTFS_DIR}/opt/splash.png"
mkdir "${ROOTFS_DIR}/var/lib/raspberryConverter"
# cp omxController.sh "${ROOTFS_DIR}/var/lib/raspberryConverter/omxController.sh"



on_chroot << EOF
# Install dependencies
apt-get update
apt-get install libpcre3 fonts-freefont-ttf fbi omxplayer -y

# Register the service
systemctl enable raspberryConverter.service

# DISABLE LOGIN
systemctl disable getty@tty1.service

# USE ALL THE SCREEN
sed -i 's/#disable_overscan=1/disable_overscan=1/g' /boot/config.txt

# SILENT BOOT
sed -i 's/ quiet splash loglevel=0 logo.nologo vt.global_cursor_default=0//g' /boot/cmdline.txt
sed -i '$ s/$/ quiet splash loglevel=0 logo.nologo vt.global_cursor_default=0/' /boot/cmdline.txt

# DISABLE RAINBOW SPLASH SCREEN
sed -i 's/disable_splash=1//g' /boot/config.txt
sed -i '$ s/$/\ndisable_splash=1/' /boot/config.txt

# CUSTOM SPLASH SCREEN
systemctl enable splashScreen.service
EOF
