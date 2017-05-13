FROM golang

# install fresh (golang rebuilder on file change)
RUN go get -u github.com/sc7639/fresh

# create app dir
RUN mkdir -p /go/src/app/ /go/src/tmp
WORKDIR /go/src/app

# Add files
ADD js/ js/
ADD template/ template/
ADD images/ images/
ADD main.go main.go

# Build app and remove source files
RUN go build && rm *.go

ENV APP_ENV prod
EXPOSE 8080

CMD ["./app"]
