FROM golang:latest as builder 
RUN mkdir /build 
ADD api/* /build/ 
WORKDIR /build
RUN go get -d -v ./
RUN go build -o main . 
RUN mkdir -p /tmp
ADD ./dyndnsConfig.json /tmp/
EXPOSE 8080
CMD ["./main"]
