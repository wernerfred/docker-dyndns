#!/bin/bash

BIND9_ROOTDOMAIN=
API_USER=
API_PASSWORD=
DYNDNS_TTL=
DYNDNS_DOMAINS=

echo "Creating local named configuration"

cat <<EOF > /etc/bind/named.conf.local
zone "${BIND9_ROOTDOMAIN}" {
    type master;
    file "/etc/bind/${BIND9_ROOTDOMAIN}.zone";
    allow-query { any; };
    allow-transfer { none; };
    allow-update { localhost; };
};
EOF

echo "Creating ${BIND9_ROOTDOMAIN} configuration"
cat <<EOF > /etc/bind/${BIND9_ROOTDOMAIN}.zone
${BIND9_ROOTDOMAIN}	IN SOA	localhost. root.localhost. (
					74         ; serial
					3600       ; refresh (1 hour)
					900        ; retry (15 minutes)
					604800     ; expire (1 week)
					86400      ; minimum (1 day)
					)
		NS	localhost.
EOF

echo "Creating named options configuration"

cat <<EOF > "/etc/bind/named.conf.options"
options {
	directory "/var/cache/bind";
    recursion no;
	dnssec-validation auto;
	auth-nxdomain no;    # conform to RFC1035
	listen-on-v6 { any; };
    allow-transfer { none; };
};
EOF

echo "Creating dyndns api config"
cat <<EOF > /tmp/
{
   "User": "${API_USER}",
   "Password": "${API_PASSWORD}",
   "Domains": "${DYNDNS_DOMAINS}",
   "TTL": "${DYNDNS_TTL}"
}
EOF


chown root:bind /var/cache/bind
chown bind:bind /var/cache/bind/*
chmod 770 /var/cache/bind
chmod 644 /var/cache/bind/*
