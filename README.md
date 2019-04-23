# raspberryConverter

***This project is in early development stage and is by no means ready to use.***

raspberryConverter is a JEOS (Just Enought Operating System) to play RTMP on Raspberry Pi 3. The project consist on two main elements: a web server that provides a web interface to configure the service and a player to show the streaming content.

## Development

### Raspberry set up
The project is meant to work under [Raspbian Stretch Lite](https://www.raspberrypi.org/downloads/raspbian/), if you are using another OS, flash and sd card with this system. If you already are using Stretch Lite, you may have to tweak some of the following steps according to your set up.

* Install git: `sudo apt-get update && sudo apt-get install git -y`
* Clone this repo and cd into it: `git clone https://github.com/AutomaticTV/raspberryConverter.git && cd raspberryConverter`
* Install dev dependencies (Go): `sh piDevEnv.sh`
* Reboot: `sudo reboot`
* Check that go is installed: `go version` should output `go version go1.12.4 linux/arm`
* Check gopath: `echo $GOPATH` should output `/home/pi/go`


### Only web set up
If you just want to test/develop the web part, the server can be run in any machine that uses docker:

* clone the project and go to the folder
* build the project: `docker build -t raspberryconverter .`
* run the server: `docker run --name raspberryconvertertest --rm -p XX:80 raspberryconverter &` where XX is a port free on your local machine.
* connect to the server: on your browser go to `http://localhost:XX`
* to stop the server once you're done: `docker kill raspberryconvertertest`
* single command of the above: `docker kill raspberryconvertertest & docker build -t raspberryconverter . && docker run --name raspberryconvertertest --rm -p 80:80 raspberryconverter &`

## The web server

The web server is using Go and Material DesignLite for the frontend.

## The player

The player is basicaly OMXPlayer controlled by the web server.
