#!/bin/sh

# Build the binary usimg docker (and move it to ./pi-gen/stage2/04-raspberry-converter/)
buildBinary(){
  docker rm -v pigen_work
  cp Dockerfile ..
  cp goDeps.sh ..
  docker build -t raspberryconvertercompiler ..
  if [ $? -ne 0 ]; then
      echo "\n\n\n\n\n\nError building the binary!"
      exit 1
  fi
  rm ../Dockerfile
  rm ../goDeps.sh
  docker run --rm -v "$@":/out raspberryconvertercompiler
  if [ $? -ne 0 ]; then
      echo "\n\n\n\n\n\nError building the binary!"
      exit 1
  fi
  echo "Binary stored at: $@"
}

update(){
  echo "==== UPDATING ===="
  echo "\n\n\n\n\n\n ==== 1. Build the binary ==== \n\n\n\n\n\n"
  buildBinary "$(pwd)"

  echo "\n\n\n\n\n\n ==== 2. Send files to Raspberry ==== \n\n\n\n\n\n"
  mkdir toSend
  cp 04-raspberry-converter/raspberryConverter.service toSend
  cp 04-raspberry-converter/splashScreen.service toSend
  cp 04-raspberry-converter/splash.png toSend
  cp 04-raspberry-converter/omxController.sh toSend
  cp raspberryConverter toSend
  echo "\n\n\n\n\n\n You must manually input the password: raspberry"
  scp -r toSend pi@"$@":/home/pi
  if [ $? -ne 0 ]; then
      rm -rf toSend
      echo "\n\n\n\n\n\nError copying files to Raspberry"
      exit 1
  fi
  rm -rf toSend

  echo "\n\n\n\n\n\n ==== 3. Install files and reboot Raspberry ==== \n\n\n\n\n\n"
  echo "You must manually input the password: raspberry"
  ssh pi@"$@" "\
  sudo systemctl stop raspberryConverter && \
  cd toSend && \
  sudo cp raspberryConverter.service /etc/systemd/system/raspberryConverter.service && \
  sudo cp raspberryConverter /etc/systemd/raspberryConverter && \
  sudo cp splashScreen.service /etc/systemd/system/splashScreen.service && \
  sudo cp splash.png /opt/splash.png && \
  sudo cp omxController.sh /var/lib/raspberryConverter/omxController.sh && \
  cd .. && rm -rf toSend && sudo reboot"
  exit 0
}

clone(){
  git clone https://github.com/RPi-Distro/pi-gen.git
  mkdir pi-gen/stage2/04-raspberry-converter
}

buildImage(){
  echo "==== BUILDING THE IMAGE ===="
  echo "\n\n\n\n\n\n ==== 1. Build the binary ==== \n\n\n\n\n\n"
  buildBinary "$(pwd)/pi-gen/stage2/04-raspberry-converter"
  cp -r 04-raspberry-converter pi-gen/stage2/
  # Build the image
  cd pi-gen
  touch stage3/SKIP
  touch stage4/SKIP
  touch stage5/SKIP
  rm stage4/EXPORT_IMAGE
  rm stage4/EXPORT_NOOBS
  rm stage5/EXPORT_IMAGE
  rm stage5/EXPORT_NOOBS

  echo "\n\n\n\n\n\n ==== 2. Build the image ==== \n\n\n\n\n\n"
  docker-compose up -d
  ./build-docker.sh
  if [ $? -ne 0 ]; then
      rm -rf toSend
      echo "\n\n\n\n\n\nError creating the image"
      cd ..
      exit 1
  fi

  echo "\n\n\n\n\n\n ==== 3. Clean ==== \n\n\n\n\n\n"
  mv deploy/*.zip ..
  cd ..
}

# Get mode
if [ "$1" = "-update" ]; then
  if [ "$2" = "" ]; then
    echo "This option requires the IP of the Raspberry to be updated."
    exit 1
  elif echo "$2" | grep -Eq '^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$'; then
    update $2
    if [ $? -ne 0 ]; then
      exit 1
    fi
    exit 0
  else
    echo "$2" is not a valid IP
    exit 1
  fi
elif [ "$1" = "-buildDev" ]
then
  clone
  cp configDev pi-gen/config
  buildImage
  exit 0
elif [ "$1" = "-buildPro" ]
then
  clone
  cp configPro pi-gen/config
  buildImage
  exit 0
else
  echo "One of the following options must be provided:\n\n\
  -update [ip address of the pi]    Update the software of the Pi based on changes made to the source.\n\
                                    The Raspberry must be accessible within the same network.\n\n\
  -buildDev                         Creates a development image (SSH enabled)\n\n\
  -buildPro                         Creates a production image (SSH disabled)"
  exit 1
fi
