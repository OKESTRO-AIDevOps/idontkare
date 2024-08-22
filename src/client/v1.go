package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	"github.com/OKESTRO-AIDevOps/idontkare/src/client/apiximpl"
	clientapiximpl "github.com/OKESTRO-AIDevOps/idontkare/src/client/apiximpl"
	"github.com/gorilla/websocket"
)

func V1SetApixImpl() {

	clientapiximpl.V1TIMEOUT_MS = CLIENT_CONFIG.TimeoutMS

}

func V1Connect(connect_url string, mani *pkgresourceapix.V1Manifest) (*websocket.Conn, error) {

	certpool := x509.NewCertPool()

	file_ca, err := os.ReadFile("cert-client/ca.pem")

	if err != nil {

		return nil, fmt.Errorf("failed to connect: read ca: %s", err.Error())
	}

	okay := certpool.AppendCertsFromPEM(file_ca)

	if !okay {

		return nil, fmt.Errorf("failed to connect: add ca: %s", err.Error())
	}

	file_cert, err := os.ReadFile("cert-client/cert.pem")

	if err != nil {
		return nil, fmt.Errorf("failed to connect: read cert: %s", err.Error())
	}

	file_cert_str := string(file_cert)

	v1mainConnectBody := make(pkgresourceapix.V1Body)

	v1mainConnectBody["cert"] = file_cert_str

	v1mainConnect, err := pkgapix.V1GetMainTemplateByAddress(pkgresourceapix.V1KindClientRequestPriv, "/connect", mani)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: failed to get connect main: %s", err.Error())

	}

	v1mainConnect.Body = v1mainConnectBody

	websocket.DefaultDialer.TLSClientConfig = &tls.Config{
		RootCAs: certpool,
	}

	c, _, err := websocket.DefaultDialer.Dial(connect_url, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: dial: %s", err.Error())
	}

	_, err = apiximpl.V1RoundTrip(v1mainConnect, c)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: connect: %s", err.Error())

	}

	return c, nil
}

func V1RunOnce(connect_url string) (string, error) {

	var result string

	if len(os.Args) < 2 {

		return "", fmt.Errorf("argc too short: %d\n", len(os.Args))

	}

	args := os.Args[1:]

	V1SetApixImpl()

	manifest, err := pkgapix.V1GetManifest()

	if err != nil {

		return "", err
	}

	v1main, err := pkgapix.V1GetMainFromArgs(pkgresourceapix.V1KindClientRequest, args, manifest)

	if err != nil {

		return "", err
	}

	conn, err := V1Connect(connect_url, manifest)

	resultData, err := clientapiximpl.V1RoundTrip(v1main, conn)

	if err != nil {

		return "", err
	}

	if resultData == nil {

		return "", fmt.Errorf("empty result data")
	}

	result = resultData.Output

	return result, nil

}

func V1Run(connect_url string, socket_addr string) error
