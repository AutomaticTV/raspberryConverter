# DEVELOPMENT

## SET UP THE DEVELOPMENT ENVIRONMENT

1. Flash a Raspbian Stretch Lite image on a SD card. (Download (here)[https://www.raspberrypi.org/downloads/raspbian/])
2. Log into your Raspberry with the flashed SD card. (user: pi, pass: raspberry)
3. Install git, clone this repo and install dependencies: `sudo apt-get install git -y && git clone https://github.com/AutomaticTV/raspberryConverter && sh raspberryConverter/dev/piDevEnv.sh`
4. Reboot: `sudo reboot`
5. Install go packages and move the project to the GO PATH (This may take a while): `sh raspberryConverter/dev/goDeps.sh && mv raspberryConverter $GOPATH/src`
6. Test that it works (Ctrl + C to stop the server, you may need to press it more than once): `cd $GOPATH/src/raspberryConverter && go run .`

Now you are ready to modify the code, and just repeat step 6 to execute.
Remember to add any go package to goDeps.sh.

## SET UP PRODUCTION LIKE ENVIRONMENT
Note that this will make changes to the boot process and you may be unable to login again unless you do so via ssh. Also, the environment here differs a bit from the development, specially with executed commands. Make sure that all required files are packed in the binary.

1. Run `cd $GOPATH/src/raspberryConverter/dev/ && sh install.sh`.
2. Reboot the system (this will launc the app on boot): `sudo reboot`.

When code changes, step 1 and 2 must be repeated to take effect.

***If for any reason install.sh is modified, the build process (creating the image for bootable SD Card) should be modified in consequence.***
