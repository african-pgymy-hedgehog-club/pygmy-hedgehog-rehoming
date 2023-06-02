FROM golang:1.18

WORKDIR /go/src/app

# Add files
ADD template/ template/
ADD *.go ./
ADD go.* ./

# Get packages
RUN go get -u github.com/sc7639/sendmail

# Build app and remove source files
RUN go build && rm *.go

ENV APP_ENV prod
EXPOSE 8080

CMD ["./app"]
