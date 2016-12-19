#!/bin/bash

set -e

rm -f ca.crt

echo "Getting latest Alpine Linux in Docker"

docker pull alpine

echo "Starting Alpine Linux in Docker and exporting certificates"

docker run --rm -v $(pwd):/roots alpine sh -c 'apk --update --no-cache add ca-certificates && touch /roots/ca.crt && cat /usr/share/ca-certificates/mozilla/*.crt >> /roots/ca.crt'

echo "Adding the certificates to cert_pool.go"

cat <<EOT > cert_pool.go
//go:generate ./update.sh

package roots

import "crypto/x509"

// MozillaRootCAs to use in API connections if x509 SystemCertPool unavailable
var MozillaRootCAs = x509.NewCertPool()

func init() {
	MozillaRootCAs.AppendCertsFromPEM([]byte(\`
EOT

cat ca.crt >> cert_pool.go

cat <<EOT >> cert_pool.go
\`))
}
EOT

echo "Done"
