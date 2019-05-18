FROM golang:1.12.4-stretch
RUN apt-get update && apt-get install -y git gcc-arm-linux-gnueabihf

WORKDIR /go/src/raspberryConverter

# Install dependencies
COPY goDeps.sh .
RUN  sh goDeps.sh && rm goDeps.sh

# Create the app
COPY frontend ./frontend
COPY models ./models
COPY services ./services
COPY main.go ./
# GOARM=6 (Raspberry Pi A, A+, B, B+, Zero) GOARM=7 (Raspberry Pi 2, 3)
RUN GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc packr build -v -o raspberryConverter
# RUN
CMD ["cp", "raspberryConverter", "/out"]
