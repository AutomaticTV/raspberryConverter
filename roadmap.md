# TO DO (first release):
### Fixtures
* DISABLE SSH
* Improve delay (research on ffmpeg / omxplayer buffers) ==> CHECK
* Remove transport ==> CHECK
* Set defaults for player settings  ==> CHECK
* Fix "fake playing" ==> omxController.sh?   ==> CHECK
### Testing
* Stable behavior on network drop  ==> CHECK
* Playing for long time

---

# Second release
* Change resolution on Display IP (not only on play)
* Refactor the code in order to get fully functional godoc
* Shutdown / Reboot button (frontend and backend)
* Keyboard sequence to set DHCP
* User input validation (it's set up but not implemented, using "gopkg.in/validator.v2" on structs)
* raspberryConverter.service specify a stop / restart mechanism
* Better "stream status". At the moment it works based on the state of the code instead of the state of the player itself.
* Check all error / success messages and improve how those are displayed on the frontend
* Check all fmt messages
* Display logs on the frontend.
* Better killing process function.
* Unify the dev and build scripts: most of the code is repeated.
* Improve the build process: clone the repo (pi-gen) instead of having the code as part of this repo.
* Use go get instead of goDeps.sh
* Speed up auth related operations: current used library takes a lot of time on ARM systems.
---

# Future releases
* FAQs / HELP frontend page
* Test suit
* Software updates
* HTTPS https://github.com/denji/golang-tls
* Make image smaller, as said on the pi-gen README: "Stage 2 - lite system. This stage produces the Raspbian-Lite image. It installs some optimized memory functions, sets timezone and charmap defaults, installs fake-hwclock and ntp, wifi and bluetooth support, dphys-swapfile, and other basics for managing the hardware. It also creates necessary groups and gives the pi user access to sudo and the standard console hardware permission groups. There are a few tools that may not make a whole lot of sense here for development purposes on a minimal system such as basic Python and Lua packages as well as the build-essential package. They are lumped right in with more essential packages presently, though they need not be with pi-gen. These are understandable for Raspbian's target audience, but ***if you were looking for something between truly minimal and Raspbian-Lite, here's where you start trimming***".
