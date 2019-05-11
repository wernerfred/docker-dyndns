# Host your own Dyndns Server using Docker

This project uses ```bind9``` and ```go``` inside a ```docker``` container to build a dyndns server/service that can easily be self-hosted. We use ```bind``` as the DNS-Server whereas ```go``` is used to serve a API and update the DNS configuration. The API uses basic authentication to restrict the usage (see [Reverse Proxy](#reverse-proxy)). You need a server with a static IP, a domain and the possibility to add ```NS``` and ```A``` records to it's DNS configuration. Furthermore it is mandatory to define the subdomains which can be used to reduce abuse in case of a data breach.

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
dyndns      IN NS      ddns
ddns        IN A       <IP of your docker host>
```
Through that we are telling the internet that the subdomain dyndns(.example.com) uses a different nameserver which can be reached under ddns(.example.com). Of course, the latter domain needs to be defined to be reached, therefore we need the second record. Keep in mind that if you want the nameserver to be a subdomain of the domain it is responsible for you will need a so called glue record on the parent zone to be able to reache the nameserver.

In case you do not like the dyndns address, you can simply add another ```CNAME``` record to beautify your dynamic address:
```
home      IN CNAME     sub1.dyndns.example.com.
````
## Using the API

### Browser
In modern browsers you can simply open the update URL:
```
https://ddns.example.com/update?domain=sub1&ip=1.2.3.4
```
You will then be promted for your API credentials that can be defined via the ```API_USER``` and ```API_PASSWORD``` variables:

![alt text](https://github.com/wernerfred/docker-dyndns/blob/master/dyndns-browser.png "Using the API via browser")
### CLI
It is also possible to use the API via command line tools like ```curl```. That command can then be used as a cronjob to constantly update the IP. The tool ```curl``` also provides the functionality to include basic authentication headers with the ```--user``` option:
```
curl --user user:password https://ddns.example.com/update?domain=sub1&ip=1.2.3.4
```
### Router
As modern routers provide a gui to configure custom dyndns services this project can also be used together with those. Usually the router uses basic authentication with the values of the user and password fields:

![alt text](https://github.com/wernerfred/docker-dyndns/blob/master/dyndns-fritzbox.png "Using the API via a router gui")

The url then looks like this:
```
https://ddns.example.com/update?domain=<domain>&ip=<ipaddr>
```
