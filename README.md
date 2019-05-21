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

## ABOUT THE CODE
All the code is documented [here](https://godoc.org/github.com/AutomaticTV/raspberryConverter). Unfortunately for the documentation to work online, a refactor is required, you can get the fully functional doc by runing `godoc -http=:6060` from the project folder. The code is mostly written in Golang, the exception is the frontend part which is implemented with HTML, CSS and JS. In order to provide a Material Design style the library [Material Design Lite](https://getmdl.io/) is used, the combination of the above makes it possible for any web frontend developer to easily modify the aspect of the frontend.

## DEVELOPMENT
This section provides a guide to set up an environment using a Raspberry Pi in order to be able to test and modify the code of this project.

### SET UP THE DEVELOPMENT ENVIRONMENT

1. Flash a Raspbian Stretch Lite image on a SD card. (Download [here](https://www.raspberrypi.org/downloads/raspbian/))
2. Log into your Raspberry with the flashed SD card. (user: pi, pass: raspberry)
3. Install git, clone this repo and install dependencies: `sudo apt-get install git -y && git clone https://github.com/AutomaticTV/raspberryConverter && sh raspberryConverter/dev/piDevEnv.sh`
4. Reboot: `sudo reboot`
5. Install go packages and move the project to the GO PATH (This may take a while): `sh raspberryConverter/dev/goDeps.sh && mv raspberryConverter $GOPATH/src`
6. Test that it works (Ctrl + C to stop the server, you may need to press it more than once): `cd $GOPATH/src/raspberryConverter && go run .`

Now you are ready to modify the code, and just repeat step 6 to execute.
Remember to add any go package to goDeps.sh.

### SET UP PRODUCTION LIKE ENVIRONMENT
Note that this will make changes to the boot process and you may be unable to login again unless you do so via ssh. Also, the environment here differs a bit from the development, specially with executed commands. Make sure that all required files are packed in the binary.

1. Run `cd $GOPATH/src/raspberryConverter/dev/ && sh install.sh`.
2. Reboot the system (this will launc the app on boot): `sudo reboot`.

When code changes, step 1 and 2 must be repeated to take effect.

***If for any reason install.sh is modified, the build process (creating the image for bootable SD Card) should be modified in consequence.***

## BUILD THE IMAGE
Unlike the development environment this process requires a x86 Linux machine (may work in macOS since it basically uses Docker, but this is not tested).

1. [Install Docker](https://docs.docker.com/install/) (if you don't have it already).
2. Install git (You should already have it).
3. Check that the following files exist: `ls /lib/modules/$(uname -r)/kernel/fs/binfmt_misc.ko` and `ls /usr/bin/qemu-arm-static`. If you don't, you have to install [binfmt_misc](https://en.wikipedia.org/wiki/Binfmt_misc)
4. (This shouldn't be necessary, but just in case) Load binfmt_misc module: `binfmt_misc`
5. Clone this repo and go to build directory: `git clone https://github.com/AutomaticTV/raspberryConverter && cd raspberryConverter/build`
6. Run the build script: `sh install.sh`
7. Do something useful with your live while install.sh creates the image
8. Once the process is done, you should have the image at build/image_YYYY-MM-DD-RaspberryConverter-lite.zip.
9. Flash the image on the SD card (you can use tools like [Etcher](https://www.balena.io/etcher/)).
10. Just put the SD card on your Raspberry.

The process of building the image has two main steps, unified by the build.sh script.

### Compile the binary
In order to produce the binary of the application, a docker image is used. This image compiles the project and outputs the binary file.

### Building the image
The project [pi-gen](https://github.com/RPi-Distro/pi-gen) is used to create the image this project is actually used by Raspbian to build their images.
A few modifications are done in order to include the required software to run the application on boot:

1. Add a file named SKIP under the folder stage3, stage4 and stage5. Remove \*Export files under this folders as well. file This will produce a image based on Rasbian lite (instead of the full desktop image).
2. Under the directory stage2 a folder will be added. The name of this folder will start with 04. This folder will contain code that will be executed in the image creation process right after the Raspbian lite is formed.
3. In the described folder an script will be placed this script will do the following: copy all the necessary files (binary, service descriptors and splash image) to the resulting root file system, register the services, install necessary packages (omxplayer, ...) and modify some files in order to produce a silent boot (just a background image with no logs).

All those changes are made using the build/build.sh script
