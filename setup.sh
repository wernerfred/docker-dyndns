#!/bin/bash

if [[ -z "${BIND9_ROOTDOMAIN}" ]];then
  echo "The variable BIND9_ROOTDOMAIN must be set"
  exit 1
fi

if [[ -z "${API_USER}" ]];then
  echo "The variable API_USER must be set"
  exit 1
fi

if [[ -z "${API_PASSWORD}" ]];then
  echo "The variable API_PASSWORD must be set"
  exit 1
fi

if [[ -z "${DYNDNS_TTL}" ]];then
  echo "The variable DYNDNS_TTL must be set"
  exit 1
fi

if [[ -z "${DYNDNS_DOMAINS}" ]];then
  echo "The variable DYNDNS_DOMAINS must be set"
  exit 1
fi

echo "Creating local named configuration"

cat <<EOF > /etc/bind/named.conf.local
zone "${BIND9_ROOTDOMAIN}" {
    type master;
    file "/var/cache/bind/${BIND9_ROOTDOMAIN}.zone";
    allow-query { any; };
    allow-transfer { none; };
    allow-update { localhost; };
};
EOF

echo "Creating ${BIND9_ROOTDOMAIN} configuration"
cat <<EOF > /var/cache/bind/${BIND9_ROOTDOMAIN}.zone

${BIND9_ROOTDOMAIN}.	IN SOA	localhost. root.localhost. (
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
cat > /root/dyndnsConfig.json <<EOF
{
   "User": "${API_USER}",
   "Password": "${API_PASSWORD}",
   "Zone": "${BIND9_ROOTDOMAIN}",
   "Domains": ${DYNDNS_DOMAINS},
   "TTL": "${DYNDNS_TTL}"
}
EOF

mkdir -p /var/cache/bind

chown root:bind /var/cache/bind
chmod 770 /var/cache/bind

service bind9 restart

/root/dyndns-api
