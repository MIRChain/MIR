package params

import "github.com/pavelkrolevets/MIR-pro/crypto/csp"


var SignerCert *csp.Cert

func SetSignerCert(cert *csp.Cert) error {
	SignerCert = cert
	return nil
}