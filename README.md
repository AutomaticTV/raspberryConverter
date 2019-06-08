# raspberryConverter
raspberryConverter is a JEOS (Just Enought Operating System) to play RTMP on Raspberry Pi 3. The project consist of two main elements: a web server that provides a web interface to configure the service and a player to show the streaming content through HDMI.

## USAGE
1. Flash the latest image provided under releases of this repository on a SD Card. You can use tools such as [Etcher](https://www.balena.io/etcher/) to do so.
2. Put the SD card on your Raspberry Pi, connect it to a monitor using the HDMI port, to the internet using the ethernet port and turn it on.
3. Wait until you see a static image with an IP address on the middle (something like http://192.168.1.35).
4. Use a web browser of any device ***that is connected to the same network of the Raspberry,*** and go to the URL you see in the monitor attached to the Raspberry (in the previous example, http://192.168.1.35).
5. Log in, the default user - password is `admin` - `admin`.
6. After the first login it is strongly recommended to change the passaword, it can be done in the PASSWORD page.
7. Next step is setting up the details of the stream you want to connect to (URL, authentication, video output resolution, ...). This is done under the PLAYER PAGE. An interesting setting is auto play, which will try to connect to the streaming server every 10 seconds, this apply even when rebooting the system.
8. Finally in the STATUS page you can start playing the video by pressing the START button, in a few seconds the video should display in the monitor connected to the Raspberry Pi.

---

## ABOUT THE CODE
All the code is documented [here](https://godoc.org/github.com/AutomaticTV/raspberryConverter). Unfortunately for the documentation to work online, a refactor is required, you can get the fully functional doc by runing `godoc -http=:6060` from the project folder.

The code is mostly written in Golang, the exception is the frontend part which is implemented with HTML, CSS and JS. In order to provide a Material Design style the library [Material Design Lite](https://getmdl.io/) is used, the combination of the above makes it possible for any web frontend developer to easily modify the aspect of the frontend.

The project [pi-gen](https://github.com/RPi-Distro/pi-gen) is used to create the image this project is actually used by Raspbian to build their images.
A few modifications are done in order to include the required software to run the application on boot:

1. Add a file named SKIP under the folder stage3, stage4 and stage5. Remove \*Export files under this folders as well. file This will produce a image based on Rasbian lite (instead of the full desktop image).
2. Under the directory stage2 a folder will be added. The name of this folder will start with 04. This folder will contain code that will be executed in the image creation process right after the Raspbian lite is formed.
3. In the described folder an script will be placed this script will do the following: copy all the necessary files (binary, service descriptors and splash image) to the resulting root file system, register the services, install necessary packages (omxplayer, ...) and modify some files in order to produce a silent boot (just a background image with no logs).

All those changes are made using the build/build.sh script.

---

## DEVELOPMENT
This section is a guide on how to modify this software and create new images. ***This process has to be run under Debian Stretch or Ubuntu Xenial (x86).***

### SET UP THE ENVIRONMENT
1. [Install Docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/) (if you don't have it already).
2. Install dependencies: `apt-get install coreutils quilt parted qemu-user-static debootstrap zerofree zip dosfstools bsdtar libcap2-bin grep rsync xz-utils file git curl`
3. Check that the following files exist: `ls /lib/modules/$(uname -r)/kernel/fs/binfmt_misc.ko` and `ls /usr/bin/qemu-arm-static`. If you don't, you have to install [binfmt_misc](https://en.wikipedia.org/wiki/Binfmt_misc)
4. (This shouldn't be necessary, but just in case) Load binfmt_misc module: `binfmt_misc`
5. Clone this repo: `git clone https://github.com/AutomaticTV/raspberryConverter`

### BUILD A DEVELOPMENT IMAGE
1. Move to the build directory (in the path of the project you've cloned): `cd ..../raspberryConverter/build`
2. Run the build script with development option: `sh build.sh -buildDev`
3. Do something useful with your live while install.sh creates the image
4. Once the process is done, you should have the image at build/image_YYYY-MM-DD-RaspberryConverterDev-lite.zip.
5. Flash the image on the SD card (you can use tools like [Etcher](https://www.balena.io/etcher/)).
6. Just put the SD card on your Raspberry.

### MODIFY THE CODE
1. Make any change to the code
2. Move to the build directory (in the path of the project you've cloned): `cd ..../raspberryConverter/build`
3. Run the build script with update option: `sh build.sh -update xxx.xxx.xxx.xxx` where xxx.xxx.xxx.xxx is the ip of the raspberry. Note that the raspberry must be in the same network as the device used to run the script. The script will ask two times for the raspberry password which is `raspberry`.
4. The raspberry will be automaticaly rebooted with the changes made on step 1.
5. To debug you can log in to the raspberry via SSH: `ssh pi@xxx.xxx.xxx.xxx` where xxx.xxx.xxx.xxx is the ip of the raspberry. Password is raspberry. The software is deployed as a systemd service. To see live logs: `sudo journalctl -f -u raspberryConverter`.

### BUILD A PRODUCTION IMAGE
Once the code changes are ready is time to build the production image, which is the same but with SSH disabled.
Follow the steps of "BUILD A DEVELOPMENT IMAGE", but this time execute the build script like this: `sh build.sh -buildPro`.
The new image will be named "image_YYYY-MM-DD-RaspberryConverter-lite.zip" instead of "image_YYYY-MM-DD-RaspberryConverterDev-lite.zip".
