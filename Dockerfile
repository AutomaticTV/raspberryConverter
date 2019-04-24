FROM golang:1.12.4-alpine
WORKDIR /go/src/raspberryConverter

# Install dependencies
COPY goDeps.sh .
RUN apk add --no-cache gcc musl-dev git && sh goDeps.sh && rm goDeps.sh

# Create the app
COPY frontend ./frontend
COPY storage ./storage
COPY services ./services
COPY main.go ./
# RUN
CMD ["go", "run", "."]
# COMPILE
# RUN CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -a -o /raspberryConverter . && chmod +x /raspberryConverter
# Run the binary
# CMD ["/raspberryConverter"]
