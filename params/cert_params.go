package params

import "github.com/MIRChain/MIR/crypto/csp"

var SignerCert *csp.Cert

func SetSignerCert(cert *csp.Cert) error {
	SignerCert = cert
	return nil
}
