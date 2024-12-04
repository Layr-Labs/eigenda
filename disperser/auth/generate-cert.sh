#!/usr/bin/env bash

if [ -z "${1}" ]; then
  echo "Usage: $0 <private-key>"
  exit 1
fi

if ! command -v cowsay 2>&1 >/dev/null
then
    echo "cowsay is not installed. Please install it ('brew install cowsay' or 'apt-get install cowsay')."
    exit 1
fi

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

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

# Document the expiration date of the certificate.
NEXT_YEAR=$("${SCRIPT_DIR}"'/next-year.py')
EXPIRATION_MESSAGE=$(cowsay "This certificate expires on ${NEXT_YEAR}.")
echo -e "${EXPIRATION_MESSAGE}\n$(cat eigenda-disperser-public.crt)" > eigenda-disperser-public.crt
cowsay "This certificate expires on ${NEXT_YEAR}. Ensure that a new one is made available to node operators before then."