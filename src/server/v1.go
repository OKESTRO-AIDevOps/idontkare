package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourcecomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/comm"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	apiximpl "github.com/OKESTRO-AIDevOps/idontkare/src/server/apiximpl"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

var V1MANI *pkgresourceapix.V1Manifest

var V1CA_CERT *x509.Certificate

func V1SetApixImpl() {

}

func V1ClientAccept(c *websocket.Conn) error {

	var req pkgresourcecomm.CommJSON
	var resp pkgresourcecomm.CommJSON

	var acceptSuccess int = 0

	var retErr error

	err := c.ReadJSON(&req)

	if err != nil {

		resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE
		resp.Message = "read failed"
		resp.Data = []byte{}

		_ = c.WriteJSON(resp)

		return fmt.Errorf("accept: failed to read: %s", err.Error())
	}

	for {

		v1mainConnect, err := pkgapix.V1GetMainByByte(req.Data, V1MANI)

		if err != nil {

			retErr = fmt.Errorf("accept: failed to get main: %s", err.Error())

			break
		}

		if v1mainConnect == nil {

			retErr = fmt.Errorf("accept: empty main: %s", err.Error())

			break
		}

		thisAddr := pkgresourceapix.V1KindClientRequestPriv + "/connect"

		targetAddr := v1mainConnect.Kind + v1mainConnect.Path

		if thisAddr != targetAddr {

			retErr = fmt.Errorf("accept: addr not matched: got: %s should: %s", targetAddr, thisAddr)

			break
		}

		certStr := v1mainConnect.Body["cert"]

		cert_b := []byte(certStr)

		clientcert, err := pkgutils.BytesToCert(cert_b)

		if err != nil {

			retErr = fmt.Errorf("accept: failed to convert: %s", err.Error())

			break
		}

		hash_sha := sha256.New()

		hash_sha.Write(clientcert.RawTBSCertificate)

		hash_data := hash_sha.Sum(nil)

		pub_key := V1CA_CERT.PublicKey.(*rsa.PublicKey)

		err = rsa.VerifyPKCS1v15(pub_key, crypto.SHA256, hash_data, clientcert.Signature)

		if err != nil {

			retErr = fmt.Errorf("accept: invalid cert: %s", err.Error())

			break
		}

		rootRecord, err := pkgdbquery.GetRoot()

		if err != nil {

			retErr = fmt.Errorf("accept: failed to get root: %s", err.Error())

			break
		}

		if rootRecord == nil {

			retErr = fmt.Errorf("accept: empty root")

			break
		}

		if rootRecord.RootName != clientcert.Subject.CommonName {

			retErr = fmt.Errorf("accept: name not matching")

			break
		}

		acceptSuccess = 1

		break

	}

	if acceptSuccess != 1 {

		resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE
		resp.Message = "not accepted"

		resp.Data = []byte{}

	} else {

		resp.Status = pkgresourcecomm.COMM_STATUS_SUCCESS
		resp.Message = "accepted"
		resp.Data = []byte{}

	}

	err = c.WriteJSON(resp)

	if err != nil {

		return fmt.Errorf("accept error: %s", err.Error())
	}

	if acceptSuccess != 1 {
		return retErr
	}

	return nil
}

func V1ClientHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("client access")

	u := websocket.Upgrader{}

	u.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("client upgrade: %s\n", err.Error())
		return
	}

	c.SetReadDeadline(time.Time{})

	defer c.Close()

	err = V1ClientAccept(c)

	if err != nil {

		log.Printf("client not accepted: %s", err.Error())

		return
	}

	log.Printf("client accepted")

	for {

		req := pkgresourcecomm.CommJSON{}

		resp := pkgresourcecomm.CommJSON{}

		err := c.ReadJSON(&req)

		if err != nil {

			log.Printf("client handle: client gone")

			_ = c.Close()

			return

		}

		v1main, err := pkgapix.V1GetMainByByte(req.Data, V1MANI)

		if err != nil {

			log.Printf("client handle: get main: %s", err.Error())

			_ = c.Close()

			return
		}

		if v1main == nil {

			log.Printf("client handle: empty main: %s", err.Error())

			_ = c.Close()

			return
		}

		var commStatus pkgresourcecomm.CommStatusType
		var commMessage string
		var v1resultReturn *pkgresourceapix.V1ResultData

		switch v1main.Kind {

		case pkgresourceapix.V1KindClientRequest:

			v1result, err := apiximpl.V1ClientRequestCtl(v1main, c)

			if err != nil {

				log.Printf("failed: clientctl: %s", err.Error())

				commStatus = pkgresourcecomm.COMM_STATUS_FAILURE

				commMessage = "failed to handle"

				v1resultReturn = nil

			} else {

				commStatus = pkgresourcecomm.COMM_STATUS_SUCCESS

				v1resultReturn = v1result

			}

		default:

			commStatus = pkgresourcecomm.COMM_STATUS_FAILURE

			commMessage = "not allowed: " + v1main.Kind

			v1resultReturn = nil

		}

		resp.Status = commStatus
		resp.Message = commMessage

		if v1resultReturn != nil {

			rr_b, err := yaml.Marshal(v1resultReturn)

			if err != nil {

				log.Printf("failed to marshal result: %s", err.Error())

				resp.Data = []byte{}

			} else {

				resp.Data = rr_b
			}

		} else {

			resp.Data = []byte{}
		}

		err = c.WriteJSON(resp)

		if err != nil {

			log.Printf("failed to write: %s", err.Error())

			_ = c.Close()

			return
		}

	}

}

func V1AgentHandler(w http.ResponseWriter, r *http.Request) {

}

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

	v1mani, err := pkgapix.V1GetManifest()

	if err != nil {

		return fmt.Errorf("failed to run: manifest: %s", err.Error())
	}

	V1MANI = v1mani

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
