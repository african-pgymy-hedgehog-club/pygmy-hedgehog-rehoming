FROM golang

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
