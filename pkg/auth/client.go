package challenge

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
)

func ClientAuthCertificate(cert_str string, ca_crt *x509.Certificate) error {

	cert_b := []byte(cert_str)

	clientcrt, err := pkgutils.BytesToCert(cert_b)

	if err != nil {

		return fmt.Errorf("client auth: invalid cert string: %s", err.Error())
	}

	hash_sha := sha256.New()

	hash_sha.Write(clientcrt.RawTBSCertificate)

	hash_data := hash_sha.Sum(nil)

	pub_key := ca_crt.PublicKey.(*rsa.PublicKey)

	err = rsa.VerifyPKCS1v15(pub_key, crypto.SHA256, hash_data, clientcrt.Signature)

	if err != nil {

		return fmt.Errorf("client auth: verify: %s", err.Error())
	}

	alleged_root := clientcrt.Subject.CommonName

	v_root, err := pkgdbquery.GetRoot()

	if err != nil {

		return fmt.Errorf("client auth: failed to get root: %s", err.Error())

	}

	if v_root.RootName != alleged_root {

		return fmt.Errorf("client auth: not root")

	}

	return nil

}
