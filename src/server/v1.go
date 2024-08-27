package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	pkgauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/auth"
	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourceauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/auth"
	pkgresourcecomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/comm"
	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	apiximpl "github.com/OKESTRO-AIDevOps/idontkare/src/server/apiximpl"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

var V1MANI *pkgresourceapix.V1Manifest

var V1CA_CERT *x509.Certificate

var agent_register = make(pkgresourceauth.AgentRegister)

var agent_address_register = make(pkgresourceauth.AgentAddressRegister)

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

			retErr = fmt.Errorf("accept: empty main")

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

func AgentAccept(c *websocket.Conn) (*pkgresourceauth.AgentRegister, error) {

	var newRegister = make(pkgresourceauth.AgentRegister)

	var req pkgresourcecomm.CommJSON
	var resp pkgresourcecomm.CommJSON

	var resultData []byte

	var retErr error

	var chalCode string

	var userRecord *pkgresourcedb.DB_User
	var clusterRecord pkgresourcedb.DB_Cluster
	var newKey string

	var acceptSuccess int = 0
	var challengeSuccess int = 0

	err := c.ReadJSON(&req)

	if err != nil {

		resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE
		resp.Message = "read failed"
		resp.Data = []byte{}

		_ = c.WriteJSON(resp)

		return nil, fmt.Errorf("accept: failed to read: %s", err.Error())
	}

	for {

		v1mainAccept, err := pkgapix.V1GetMainByByte(req.Data, V1MANI)

		if err != nil {

			retErr = fmt.Errorf("accept: failed to get main: %s", err.Error())

			break
		}

		if v1mainAccept == nil {

			retErr = fmt.Errorf("accept: empty main")

			break
		}

		thisAddr := pkgresourceapix.V1KindAgentRequestPriv + "/cluster/connect"

		targetAddr := v1mainAccept.Kind + v1mainAccept.Path

		if thisAddr != targetAddr {

			retErr = fmt.Errorf("accept: addr not matched: got: %s should: %s", targetAddr, thisAddr)

			break
		}

		thisClusterName, okay := v1mainAccept.Body["name"]

		if !okay {

			retErr = fmt.Errorf("accept: 'name' doesn't exist")

			break
		}

		thisUserName, okay := v1mainAccept.Body["username"]

		if !okay {

			retErr = fmt.Errorf("accept: 'username' doesn't exist")

			break
		}

		urecord, err := pkgdbquery.GetUserByName(thisUserName)

		if err != nil {

			retErr = fmt.Errorf("accept: failed to get user: %s", thisUserName)

			break
		}

		if urecord == nil {

			retErr = fmt.Errorf("accept: empty user record")

			break
		}

		userRecord = urecord

		crecord, err := pkgdbquery.GetClustersByUserId(urecord.UserId)

		if err != nil {

			retErr = fmt.Errorf("accept: failed to cluster record: %s", err.Error())

			break
		}

		clen := len(crecord)

		if crecord == nil || clen == 0 {

			retErr = fmt.Errorf("accept: cluster zero")

			break

		}

		c_found := -1

		for i := 0; i < clen; i++ {

			if crecord[i].ClusterName == thisClusterName {
				c_found = i
				break
			}

		}

		if c_found == -1 {

			retErr = fmt.Errorf("accept: cluster not found: %s", err.Error())

			break
		}

		clusterRecord = crecord[c_found]

		cpub_str := crecord[c_found].ClusterPub

		cpub, err := pkgutils.BytesToPublicKey([]byte(cpub_str))

		if err != nil {

			retErr = fmt.Errorf("accept: get pubkey: %s", err.Error())
			break
		}

		thisChalCode, _ := pkgutils.RandomHex(32)

		chalCode = thisChalCode

		thisChalCode_b := []byte(thisChalCode)

		thisChalCode_enc, err := pkgutils.EncryptWithPublicKey(thisChalCode_b, cpub)

		thisChalCode_hex := hex.EncodeToString(thisChalCode_enc)

		resp_output := pkgresourceapix.V1ResultData{}

		resp_output.Output = thisChalCode_hex

		result_yb, err := yaml.Marshal(resp_output)

		if err != nil {

			retErr = fmt.Errorf("accept: marshal result: %s", err.Error())

			break
		}

		resultData = result_yb

		acceptSuccess = 1

		break

	}

	if acceptSuccess != 1 {

		log.Printf("failed accept: %s\n", retErr.Error())

		resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE
		resp.Message = "not accepted"

		resp.Data = []byte{}
	} else {

		resp.Status = pkgresourcecomm.COMM_STATUS_SUCCESS
		resp.Message = "challenge accepted"

		resp.Data = resultData
	}

	err = c.WriteJSON(resp)

	if err != nil {

		return nil, fmt.Errorf("accept error: %s", err.Error())
	}

	if acceptSuccess != 1 {
		return nil, retErr
	}

	req = pkgresourcecomm.CommJSON{}
	resp = pkgresourcecomm.CommJSON{}

	resultData = []byte{}

	retErr = nil

	err = c.ReadJSON(&req)

	if err != nil {

		resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE
		resp.Message = "read failed"
		resp.Data = []byte{}

		_ = c.WriteJSON(resp)

		return nil, fmt.Errorf("accept: failed to read: %s", err.Error())
	}

	// TODO:
	//   needs timeout

	for {

		v1mainChal, err := pkgapix.V1GetMainByByte(req.Data, V1MANI)

		if err != nil {

			retErr = fmt.Errorf("chal: failed to get main: %s", err.Error())

			break
		}

		if v1mainChal == nil {

			retErr = fmt.Errorf("chal: empty main")

			break
		}

		thisAddr := pkgresourceapix.V1KindAgentRequestPriv + "/cluster/connect/challenge"

		targetAddr := v1mainChal.Kind + v1mainChal.Path

		if thisAddr != targetAddr {

			retErr = fmt.Errorf("chal: addr not matched: got: %s should: %s", targetAddr, thisAddr)

			break
		}

		chal_data, okay := v1mainChal.Body["chaldata"]

		if !okay {

			retErr = fmt.Errorf("chal: 'chaldata' doesn't exist: %s", err.Error())

			break
		}

		chalCode_b := []byte(chalCode)

		chal_data_b, err := hex.DecodeString(chal_data)

		if err != nil {

			retErr = fmt.Errorf("chal: decode: %s", err.Error())

			break
		}

		data_b, err := pkgutils.DecryptWithSymmetricKey(chalCode_b, chal_data_b)

		if err != nil {

			retErr = fmt.Errorf("chal: decrypt: %s", err.Error())

			break
		}

		cdata := pkgresourceauth.ChallengeData{}

		err = json.Unmarshal(data_b, &cdata)

		if err != nil {

			retErr = fmt.Errorf("chal: unmarshal chaldata: %s", err.Error())

			break
		}

		if userRecord.UserPass != cdata.Pass {

			retErr = fmt.Errorf("chal: pass not matched: %s", err.Error())

			break
		}

		err = pkgauth.VerifySessionKey(cdata.Key)

		if err != nil {

			retErr = fmt.Errorf("chal: key invalid: %s", err.Error())

			break
		}

		newKey = cdata.Key

		newKey_b := []byte(newKey)

		tmp_message := "okay"

		enc_b, err := pkgutils.EncryptWithSymmetricKey(newKey_b, []byte(tmp_message))

		if err != nil {

			retErr = fmt.Errorf("chal: encryp with new: %s", err.Error())

			break

		}

		enc_hex := hex.EncodeToString(enc_b)

		result_data := pkgresourceapix.V1ResultData{}

		result_data.Output = enc_hex

		yb, err := yaml.Marshal(result_data)

		if err != nil {

			retErr = fmt.Errorf("chal: marshal result: %s", err.Error())

			break
		}

		resultData = yb

		challengeSuccess = 1

		break

	}

	if challengeSuccess != 1 {

		log.Printf("failed chal: %s\n", retErr.Error())

		resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE
		resp.Message = "challenge failed"

		resp.Data = []byte{}

	} else {

		resp.Status = pkgresourcecomm.COMM_STATUS_SUCCESS
		resp.Message = "cahllenge success"

		resp.Data = resultData

	}

	err = c.WriteJSON(resp)

	if err != nil {

		return nil, fmt.Errorf("chal error: %s", err.Error())
	}

	if challengeSuccess != 1 {
		return nil, retErr
	}

	registerKey := userRecord.UserName + ":" + clusterRecord.ClusterName

	agentData := pkgresourceauth.AgentData{
		Key: newKey,
		C:   c,
	}

	newRegister[registerKey] = agentData

	return &newRegister, nil
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

	log.Println("agent access")

	u := websocket.Upgrader{}

	u.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("agent upgrade: %s\n", err.Error())
		return
	}

	c.SetReadDeadline(time.Time{})

	defer c.Close()

	new_ar, err := AgentAccept(c)

	if err != nil {

		log.Printf("agent accept: %s\n", err.Error())

		return
	}

	thisAgentId := ""

	for k, v := range *new_ar {

		agent_register[k] = v

		thisAgentId = k
	}

	agent_address_register[c] = thisAgentId

	log.Printf("agent accepted")

	keep_running := 1

	for keep_running == 1 {

	}
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
