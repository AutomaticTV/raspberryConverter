# Build

This README will walk you through the process of creating a OS that consist of a clean version of Rasbian Stretch Lite and the code of this project.

Building the image takes something like 10 - 30 minutes, so it's recommended to firs try any change in the development environtment (see dev/README.md).

## SET UP PRODUCTION ENVIRONMENT
Unlike the development environment this process requires a x86 Linux machine.

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
