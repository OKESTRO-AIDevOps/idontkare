package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
)

var V1CA_CERT *x509.Certificate

func V1SetApixImpl() {

}

func V1ClientHandler(w http.ResponseWriter, r *http.Request)

func V1AgentHandler(w http.ResponseWriter, r *http.Request)

func V1Run() error {

	err := pkgdbquery.DbEstablish(
		SERVER_CONFIG.DBId,
		SERVER_CONFIG.DBPw,
		SERVER_CONFIG.DBAddr,
		SERVER_CONFIG.DBName,
	)

	if err != nil {

		return fmt.Errorf("failed to run: db: %s", err.Error())

	}

	log.Printf("db connected: %s", SERVER_CONFIG.DBAddr)

	file_b, err := os.ReadFile("cert-server/ca.pem")

	if err != nil {

		return fmt.Errorf("failed to run: ca file: %s", err.Error())
	}

	cacert, err := pkgutils.BytesToCert(file_b)

	if err != nil {

		return fmt.Errorf("failed to run: make ca: %s", err.Error())
	}

	V1CA_CERT = cacert

	listen_addr_full := SERVER_CONFIG.ListenAddr + ":" + SERVER_CONFIG.ListenPort

	http.HandleFunc(SERVER_CONFIG.ClientPath, V1ClientHandler)
	http.HandleFunc(SERVER_CONFIG.AgentPath, V1AgentHandler)

	log.Printf("server started at: %s", listen_addr_full)

	log.Fatal(http.ListenAndServeTLS(listen_addr_full, "cert-server/cert.pem", "cert-server/cert-priv.pem", nil))

	return nil
}
