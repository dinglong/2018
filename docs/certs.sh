#!/bin/bash

function gen_ca() {
  local ca=$1
  local name=$2

  openssl genrsa -out ${name}.key 2048
  openssl req -subj /CN=${name} -new -key ${name}.key -out ${name}.csr

  export v3_ca_ext="[v3_ca]\n
               subjectKeyIdentifier = hash\n
               authorityKeyIdentifier = keyid:always,issuer\n
               basicConstraints = CA:true"

  openssl ca -batch -create_serial -days 5000 -policy policy_anything -extensions v3_ca -extfile <(echo -e ${v3_ca_ext}) -keyfile ${ca}.key -cert ${ca}.crt -in ${name}.csr -out ${name}.crt
}

function gen_cert() {
  local ca=$1
  local name=$2
  local subject=$3
  local alt_name=$4

  openssl genrsa -out ${name}.key 2048

  config="[req]\n
    distinguished_name=req_distinguished_name\n
  [req_distinguished_name]\n
  [SAN]\n
  subjectAltName=${alt_name}\n"

  openssl req -new -subj ${subject} -key ${name}.key -out ${name}.csr -reqexts SAN -config <(echo -e "${config}")

  export v3_req_ext="[v3_req]\n
    subjectAltName=${alt_name}"
  
  openssl ca -batch -policy policy_anything -days 5000 -extensions v3_req -extfile <(echo -e "${v3_req_ext}") -keyfile ${ca}.key -cert ${ca}.crt -in ${name}.csr -out ${name}.cert
}

if [ ! -f "/etc/pki/CA/index.txt" ]; then
  touch /etc/pki/CA/index.txt
  echo '1000' > /etc/pki/CA/serial
fi

# rootca
openssl genrsa -out rootca.key 2048
openssl req -new -subj /CN=rootca -key rootca.key -out rootca.csr
openssl ca -batch -selfsign -create_serial -days 5000 -policy policy_anything -extensions v3_ca -keyfile rootca.key -in rootca.csr -out rootca.crt

# middle ca
gen_ca rootca ca-middle

# server cert
gen_cert ca-middle server /CN=server DNS:localhost,DNS:server,IP:127.0.0.1

# client cert
gen_cert ca-middle client /CN=client DNS:localhost,DNS:client,IP:127.0.0.1
