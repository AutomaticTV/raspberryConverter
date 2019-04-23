# raspberryConverter

*This project is in early development stage and is by no means ready to use.*

raspberryConverter is a JEOS (Just Enought Operating System) to play RTMP on Raspberry Pi 3. The project consist on two main elements: a web server that provides a web interface to configure the service and a player to show the streaming content.

## Usage

If you just want to test/develop the web part, the server can be run in any machine that uses docker:
* clone the project and go to the folder
* build the project: `docker build -t rtmpi .`
* run the server: `docker run --name rtmpitest --rm -p XX:80 &` where XX is a port free on your local machine.
* connect to the server: on your browser go to `http://localhost:XX`
* to stop the server once you're done: `docker kill rtmpitest`

## The web server

The web server is using Go and Material DesignLite for the frontend.

## The player

The player is basicaly OMXPlayer controlled by the web server.
