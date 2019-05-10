# Host your own Dyndns Server using Docker

## Installation

### Build from source
To build this project from source make sure to clone the repository from github and change your working directory into that newly created folder. Then simply issue the ```docker build``` command like shown:
```
root@server /opt/docker-dyndns # docker build -t wernerfred/docker-dyndns .
```
To run the container adjust the following command according to your needs:
```
docker run -it -d \
           -p 53:53 \ 
           -p 53:53/udp \
           -p 8080:8080 \
           -e BIND9_ROOTDOMAIN=dyndns.example.com \
           -e API_USER=user \
           -e API_PASSWORD=password \
           -e DYNDNS_TTL=60 \
           -e DYNDNS_DOMAINS='["sub1", "sub2"]' \
           wernerfred/docker-dyndns
```
With the variable ```BIND9_ROOTDOMAIN``` you set the domain of your dyndns server (See [DNS configuration](#dns-configuration)). With ```DYNDNS_DOMAINS``` you define the allowed dynamic subdomains in a json like array. Subdomains not listed there can not be updated later on. This keeps the consequences relatively small in case your API username and password gets exposed as only the defined subdomains can be updated.
### Reverse proxy
I highly recommend to use a reverse proxy on your system to kind of secure the API access. If you are using ```apache2``` you can orient yourself on the following configuration:
```
<IfModule mod_ssl.c>
<VirtualHost *:443>

    ProxyPreserveHost On
    ProxyRequests off
    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/

    ServerName ddns.example.com

    ServerAdmin webmaster@example.com

    ErrorLog /your/path/to/logs/error.log
    CustomLog /your/path/to/logs/access.log combined

SSLCertificateFile /your/path/to/fullchain.pem
SSLCertificateKeyFile /your/path/to/privkey.pem

</VirtualHost>
</IfModule>
```
Make sure to include a redirect to https (:443) into the http (:80) configuration file like that:
```
RewriteEngine on
RewriteCond %{SERVER_NAME} =ddns.example.com
RewriteRule ^ https://%{SERVER_NAME}%{REQUEST_URI} [END,NE,R=permanent]
```
### DNS configuration
To allow your dyndns server to be reached and used you need to add some DNS records of your existing domain. Like in that whole README i assume you own the domain ```example.com```. If the dyndns domains should use ```dyndns.example.com``` as their root domain we need the following records:
```
dyndns                   IN NS      ddns
ddns                     IN A       <IP of your docker host>
```
Through that we are telling the internet that the subdomain dyndns(.example.com) uses a different nameserver which can be reached under ddns(.example.com). Of course, the latter domain needs to be defined to be reached, therefore we need the second record. Keep in mind that if you want the nameserver to be a subdomain of the domain it is responsible for you will need a so called glue record on the parent zone to be able to reache the nameserver.
## Using the API

### Browser

### CLI

### Router
