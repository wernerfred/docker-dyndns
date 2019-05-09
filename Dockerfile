FROM golang:latest as builder 
RUN mkdir /build 
ADD api/* /build/ 
WORKDIR /build
RUN go get -d -v ./
RUN go build -o dyndns-api . 

FROM debian:stretch
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update -qq && \
    apt install -q -y bind9 dnsutils -qq && \
    apt-get clean

COPY --from=builder /build/dyndns-api /root/


ADD ./dyndnsConfig.json /tmp/
EXPOSE 8080
CMD ["/root/dyndns-api"]
