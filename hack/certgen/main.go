package main

import (
	"fmt"
	"os"

	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
)

func ToFile(cs *pkgutils.CertSet) error {

	err := os.WriteFile("cert-server/ca.pem", cs.RootCertPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())

	}

	err = os.WriteFile("cert-server/priv.pem", cs.RootKeyPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-server/pub.pem", cs.RootPubPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-client/ca.pem", cs.RootCertPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())

	}

	err = os.WriteFile("cert-client/priv.pem", cs.RootKeyPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-client/pub.pem", cs.RootPubPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-server/cert.pem", cs.ServCertPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-server/cert-priv.pem", cs.ServKeyPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-server/cert-pub.pem", cs.ServPubPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-client/cert.pem", cs.ClientCertPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-client/cert-priv.pem", cs.ClientKeyPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	err = os.WriteFile("cert-client/cert-pub.pem", cs.ClientPubPEM, 0644)

	if err != nil {

		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	return nil

}

func main() {

	arglen := len(os.Args)

	if arglen != 3 {

		fmt.Fprintf(os.Stderr, "arglen invalid: %d\n", arglen)

		return
	}

	hostName := os.Args[1]
	rootName := os.Args[2]

	var cs *pkgutils.CertSet

	cs = pkgutils.NewCertsPipeline(hostName, rootName)

	err := ToFile(cs)

	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err.Error())

	} else {

		fmt.Fprintf(os.Stdout, "certs successfully generated for: host: %s, root: %s\n", hostName, rootName)

	}

	return
}
