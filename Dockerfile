FROM golang

# Install sendmail
RUN apt-get update && apt-get install sendmail -y
# Set FQDN
RUN line=$(head -n 1 /etc/hosts) && line2=$(echo $line | awk '{print $2}') && echo "$line $line2.localdomain" >> /etc/hosts

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
