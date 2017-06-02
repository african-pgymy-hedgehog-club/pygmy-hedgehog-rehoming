FROM golang

WORKDIR /go/src/app

# Get packages
RUN go get -u github.com/sc7639/sendmail

# Add files
ADD js/ js/
ADD template/ template/
ADD images/ images/
ADD *.go ./

# Build app and remove source files
RUN go build && rm *.go

ENV APP_ENV prod
EXPOSE 8080

CMD ["./app"]
