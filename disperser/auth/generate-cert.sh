#!/usr/bin/env bash

if [ -z "${1}" ]; then
  echo "Usage: $0 <private-key>"
  exit 1
fi

# Used to pass options to the 'openssl req' command.
# It expects a human in the loop, but it's preferable to automate it.
options() {
  # Country Name (2 letter code) [AU]:
  echo 'US'
  # State or Province Name (full name) [Some-State]:
  echo 'Washington'
  # Locality Name (eg, city) []:
  echo 'Seattle'
  # Organization Name (eg, company) [Internet Widgits Pty Ltd]:
  echo 'Eigen Labs'
  # Organizational Unit Name (eg, section) []:
  echo 'EigenDA'
  # Common Name (e.g. server FQDN or YOUR name) []:
  echo 'disperser'
  # Email Address []:
  echo '.'
  # A challenge password []:
  echo '.'
  # An optional company name []:
  echo '.'
}

# Generate a new certificate signing request.
options | \
  openssl req -new \
  -key "${1}" \
  -noenc \
  -out cert.csr

# Self sign the certificate.
openssl x509 -req \
  -days 365 \
  -in cert.csr \
  -signkey "${1}" -out eigenda-disperser-public.crt

# Clean up the certificate signing request.
rm cert.csr

cowsay "This certificate will expire in one year. Ensure that a new one is made available to node operators before then."