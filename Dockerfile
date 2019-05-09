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
COPY ./setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh 

EXPOSE 53 8080
ENTRYPOINT ["/root/setup.sh"]
