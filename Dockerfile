FROM golang:1.12.4-alpine
WORKDIR /go/src/app
# Install dependencies
RUN apk add --no-cache gcc musl-dev git && go get \
github.com/gorilla/sessions \
golang.org/x/crypto/bcrypt \
github.com/jinzhu/gorm \
github.com/jinzhu/gorm/dialects/sqlite
# TODO: run a script with 3 options:
#       serve: the current one. (default)
#       test: ?
#       pack: create binary using https://github.com/gobuffalo/packr => export to a release folder (docker volume)
COPY src .
CMD ["go", "run", "."]
