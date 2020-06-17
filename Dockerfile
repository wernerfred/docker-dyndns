FROM golang:1.14.1 as builder 
RUN mkdir /build 
COPY api/* /build/ 
WORKDIR /build
RUN go get -d -v ./
RUN go build -o dyndns-api . 

FROM debian:stretch
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -qq && \
    apt-get install -q -y bind9=1:9.10.3.dfsg.P4-12.3+deb9u6 dnsutils=1:9.10.3.dfsg.P4-12.3+deb9u6 -qq --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/dyndns-api /root/
COPY ./setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh 

EXPOSE 53 8080
ENTRYPOINT ["/root/setup.sh"]