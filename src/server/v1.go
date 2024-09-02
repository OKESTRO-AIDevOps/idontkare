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
	"strings"
	"time"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	pkgauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/auth"
	pkgcomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/comm"
	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourceauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/auth"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
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

func V1AgentAccept(c *websocket.Conn) (*pkgresourceauth.AgentRegister, error) {

	var newRegister = make(pkgresourceauth.AgentRegister)

	var req pkgresourcecomm.CommJSON
	var resp pkgresourcecomm.CommJSON

	var resultData []byte

	var retErr error

	var chalCode string

	var userName string

	var userRecord pkgresourcedb.DB_User
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

		userRecord = *urecord

		userName = thisUserName

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

		thisChalCode, _ := pkgutils.RandomHex(16)

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

			retErr = fmt.Errorf("chal: pass not matched")

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
		resp.Message = "challenge success"

		resp.Data = resultData

	}

	err = c.WriteJSON(resp)

	if err != nil {

		return nil, fmt.Errorf("chal error: %s", err.Error())
	}

	if challengeSuccess != 1 {
		return nil, retErr
	}

	registerKey := userName + ":" + clusterRecord.ClusterName

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

			v1result, err := apiximpl.V1ClientRequestCtl(v1main)

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

	new_ar, err := V1AgentAccept(c)

	if err != nil {

		log.Printf("agent accept: %s\n", err.Error())

		return
	}

	thisAgentName := ""
	thisAgentCluster := ""

	thisAgentId, thisAgentData, err := V1AgentAdd(new_ar)

	if err != nil {
		log.Printf("agent accept: add: %s\n", err.Error())

		return
	}

	this_agent_li := strings.SplitN(thisAgentId, ":", 2)

	thisAgentName = this_agent_li[0]
	thisAgentCluster = this_agent_li[1]

	log.Printf("========== agent accepted")

	fmt.Println(agent_register)

	log.Printf("==========")

	for {

		req := pkgresourcecomm.CommJSON{}

		// for now, only push exists
		// not sending response
		// resp := pkgresourcecomm.CommJSON{}

		err := c.ReadJSON(&req)

		if err != nil {

			delerr := V1AgentDel(thisAgentId)

			if delerr != nil {

				log.Printf("agent handle: agent gone: failed to delete: %s: %s\n", thisAgentId, delerr.Error())

			} else {

				log.Printf("agent handle: agent gone: successfully deleted: %s\n", thisAgentId)
			}

			_ = c.Close()

			return
		}

		data_b, err := pkgcomm.CommDataDecrypt(string(req.Data), []byte(thisAgentData.Key))

		if err != nil {

			log.Printf("id: %s: decrypt: %s\n", thisAgentId, err.Error())

			continue
		}

		v1main, err := pkgapix.V1GetMainByByte(data_b, V1MANI)

		if err != nil {
			log.Printf("id: %s: main: %s\n", thisAgentId, err.Error())

			continue
		}

		if v1main == nil {

			log.Printf("id: %s: empty\n", thisAgentId)

			continue
		}

		err = apiximpl.V1AgentPush(v1main, thisAgentName, thisAgentCluster)

		if err != nil {

			log.Printf("id: %s: agent push: %s", thisAgentId, err.Error())
		}

	}
}

func V1AgentAdd(new_ar *pkgresourceauth.AgentRegister) (string, *pkgresourceauth.AgentData, error) {

	thisAgentId := ""
	thisAgentData := pkgresourceauth.AgentData{}

	for k, v := range *new_ar {

		agent_register[k] = v

		thisAgentData = v

		thisAgentId = k

	}

	this_agent_li := strings.SplitN(thisAgentId, ":", 2)

	thisAgentName := this_agent_li[0]
	thisAgentCluster := this_agent_li[1]

	urecord, err := pkgdbquery.GetUserByName(thisAgentName)

	if err != nil {

		delete(agent_register, thisAgentId)

		return "", nil, fmt.Errorf("failed to add agent: %s", err.Error())
	}

	if urecord == nil {

		delete(agent_register, thisAgentId)

		return "", nil, fmt.Errorf("failed to add agent: empty record")
	}

	err = pkgdbquery.SetClusterConnectedByUserIdAndName(urecord.UserId, thisAgentCluster, 1, thisAgentData.Key)

	if err != nil {

		delete(agent_register, thisAgentId)

		return "", nil, fmt.Errorf("failed to add agent: set cluster: %s", err.Error())
	}

	return thisAgentId, &thisAgentData, nil
}

func V1AgentDel(agent_id string) error {

	this_agent_li := strings.SplitN(agent_id, ":", 2)

	thisAgentName := this_agent_li[0]
	thisAgentCluster := this_agent_li[1]

	urecord, err := pkgdbquery.GetUserByName(thisAgentName)

	if err != nil {

		return fmt.Errorf("failed to del agent: %s", err.Error())
	}

	if urecord == nil {

		return fmt.Errorf("failed to del agent: empty record")
	}

	err = pkgdbquery.SetClusterConnectedByUserIdAndName(urecord.UserId, thisAgentCluster, 0, "")

	if err != nil {

		return fmt.Errorf("failed to del agent: set cluster: %s", err.Error())
	}

	delete(agent_register, agent_id)

	return nil

}

func V1ProjectControlLoop() {

	var FAIL_COUNT_LIMIT int = 100

	fail_count := 0

	for {

		var CI_Q = make([]pkgresourceci.CiOption, 0)
		var CD_Q = make([]pkgresourcecd.CdOption, 0)

		if fail_count >= FAIL_COUNT_LIMIT {

			log.Fatalf("pctl: fail_count limit exceeded: %d\n", fail_count)
		}

		// poll project ci cd

		PROJECTS, err := pkgdbquery.GetProject()

		if err != nil {

			log.Printf("pctl: failed to get project: %s\n", err.Error())

			fail_count += 1

			continue
		}

		PLEN := len(PROJECTS)

		for i := 0; i < PLEN; i++ {

			//projectid := PROJECTS[i].ProjectId
			//userid := PROJECTS[i].UserId

			cdopt := pkgresourcecd.CdOption{}
			ciopt := pkgresourceci.CiOption{}

			var cdopt_err error
			var ciopt_err error

			var process_err error = nil

			cdoption_raw := PROJECTS[i].ProjectCdOption

			cioption_raw := PROJECTS[i].ProjectCiOption

			if !cdoption_raw.Valid && !cioption_raw.Valid {

				//log.Printf("pctl: skip: both null: uid: %d, pid: %d", userid, projectid)

				continue

			} else if !cdoption_raw.Valid && cioption_raw.Valid {

				ciopt_err = yaml.Unmarshal([]byte(cioption_raw.String), &ciopt)

				if ciopt.Request == nil {

					// log.Printf("pctl: skip: ci: already processed: uid: %d, pid: %d", userid, projectid)

					continue
				}

				if ciopt_err != nil {

					process_err = fmt.Errorf("pctl: skip: ci unmarshal: %s", ciopt_err.Error())

				}

				if process_err != nil {

					ciopt.Process = &struct {
						StoredRequest pkgresourceci.CiOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LinkToCd      int                            "yaml:\"link_to_cd\""
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *ciopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LinkToCd:      -1,
						Error:         true,
						Log:           process_err.Error(),
					}

				} else {
					ciopt.Process = &struct {
						StoredRequest pkgresourceci.CiOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LinkToCd      int                            "yaml:\"link_to_cd\""
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *ciopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LinkToCd:      -1,
						Error:         false,
						Log:           "",
					}

				}

				ciopt.Request = nil

				CI_Q = append(CI_Q, ciopt)

			} else if cdoption_raw.Valid && !cioption_raw.Valid {

				cdopt_err = yaml.Unmarshal([]byte(cdoption_raw.String), &cdopt)

				if cdopt.Request == nil {

					log.Printf("pctl: skip: cd: already processed")

					continue
				}

				if cdopt_err != nil {

					process_err = fmt.Errorf("pctl: skip: cd unmarshal: %s", cdopt_err.Error())

				}

				if cdopt.Request.DependOnCI {

					if process_err != nil {

						process_err = fmt.Errorf("pctl: depend on ci, but null ci + %s  ", process_err.Error())

					} else {

						process_err = fmt.Errorf("pctl: depend on ci, but null ci")

					}

				}

				if process_err != nil {

					cdopt.Process = &struct {
						StoredRequest pkgresourcecd.CdOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LifecycleId   int                            `yaml:"lifecycle_id"`
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *cdopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LifecycleId:   -1,
						Error:         true,
						Log:           process_err.Error(),
					}

				} else {
					cdopt.Process = &struct {
						StoredRequest pkgresourcecd.CdOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LifecycleId   int                            `yaml:"lifecycle_id"`
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *cdopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LifecycleId:   -1,
						Error:         false,
						Log:           "",
					}

				}

				cdopt.Request = nil

				CD_Q = append(CD_Q, cdopt)

			} else if cdoption_raw.Valid && cioption_raw.Valid {

				should_process_ci := 1
				should_process_cd := 1

				ciopt_err = yaml.Unmarshal([]byte(cioption_raw.String), &ciopt)

				cdopt_err = yaml.Unmarshal([]byte(cdoption_raw.String), &cdopt)

				if ciopt_err != nil {

					ciopt_err = fmt.Errorf("ci failed to unmarshal: %s", ciopt_err.Error())
				}

				if cdopt_err != nil {

					cdopt_err = fmt.Errorf("cd failed to unamrshal: %s", cdopt_err.Error())

				}

				if ciopt.Request == nil {

					should_process_ci = 0
				}

				if cdopt.Request == nil {

					should_process_cd = 0
				}

				if should_process_cd == 0 && should_process_ci == 0 {

					//log.Printf("pctl: skip: ci cd: already processed")

					continue
				}

				if ciopt_err != nil && should_process_ci == 1 {

					ciopt.Process = &struct {
						StoredRequest pkgresourceci.CiOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LinkToCd      int                            "yaml:\"link_to_cd\""
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *ciopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LinkToCd:      -1,
						Error:         true,
						Log:           ciopt_err.Error(),
					}
				} else if ciopt_err == nil && should_process_ci == 1 {

					ciopt.Process = &struct {
						StoredRequest pkgresourceci.CiOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LinkToCd      int                            "yaml:\"link_to_cd\""
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *ciopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LinkToCd:      -1,
						Error:         false,
						Log:           "",
					}
				}

				if cdopt_err != nil && should_process_cd == 1 {

					cdopt.Process = &struct {
						StoredRequest pkgresourcecd.CdOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LifecycleId   int                            `yaml:"lifecycle_id"`
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *cdopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LifecycleId:   -1,
						Error:         true,
						Log:           cdopt_err.Error(),
					}
				} else if cdopt_err == nil && should_process_cd == 1 {

					cdopt.Process = &struct {
						StoredRequest pkgresourcecd.CdOption_Request "yaml:\"stored_request\""
						ProjectIndex  int                            "yaml:\"project_index\""
						ProjectId     int                            `yaml:"project_id"`
						UserId        int                            `yaml:"user_id"`
						ProjectName   string                         `yaml:"project_name"`
						LifecycleId   int                            `yaml:"lifecycle_id"`
						Error         bool                           `yaml:"error"`
						Log           string                         `yaml:"log"`
					}{
						StoredRequest: *cdopt.Request,
						ProjectIndex:  i,
						ProjectId:     PROJECTS[i].ProjectId,
						UserId:        PROJECTS[i].UserId,
						ProjectName:   PROJECTS[i].ProjectName,
						LifecycleId:   -1,
						Error:         false,
						Log:           "",
					}
				}

				ciopt.Request = nil
				cdopt.Request = nil

				if cdopt.Process.StoredRequest.DependOnCI {

					linkidx := len(CD_Q)

					ciopt.Process.LinkToCd = linkidx

				}

				if should_process_ci == 1 {

					CI_Q = append(CI_Q, ciopt)
				}

				if should_process_cd == 1 {

					CD_Q = append(CD_Q, cdopt)

				}

			}
		}

		CI_Q_LEN := len(CI_Q)
		CD_Q_LEN := len(CD_Q)

		// push ci

		for i := 0; i < CI_Q_LEN; i++ {

			cioption := CI_Q[i]

			if cioption.Process.Error {

				err = V1PCTL_HandleCiError(&cioption, nil)

				if err != nil {

					fail_count += 1

					log.Printf("pctl: skip: error while processing ci error: %s", err.Error())

				}

				continue
			}

			user_clusters, err := pkgdbquery.GetClustersByUserId(cioption.Process.UserId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while querying user cluster")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			user_projectcis, err := pkgdbquery.GetProjectCisByProjectId(PROJECTS[cioption.Process.ProjectIndex].ProjectId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while querying user project cis")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			electedCluster, err := V1PCTL_ElectCiAllocableClusterId(user_clusters, user_projectcis)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while electing cluster for ci")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			newpci, err := pkgdbquery.SetProjectCi(PROJECTS[cioption.Process.ProjectIndex].ProjectId, electedCluster.ClusterId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while setting project ci")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			v1main, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindServerPush, "/project/cluster/ci/alloc", V1MANI)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while getting main")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				_ = pkgdbquery.SetProjectCiEndById(newpci.ProjectCiId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			yb, err := yaml.Marshal(cioption)

			if err != nil {
				fail_count += 1

				this_error := fmt.Errorf("error while preparing main")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				_ = pkgdbquery.SetProjectCiEndById(newpci.ProjectCiId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			pidx := cioption.Process.ProjectIndex

			v1main.Body["name"] = cioption.Process.ProjectName
			v1main.Body["git"] = PROJECTS[pidx].ProjectGit
			v1main.Body["gitid"] = PROJECTS[pidx].ProjectGitId
			v1main.Body["gitpw"] = PROJECTS[pidx].ProjectGitPw
			v1main.Body["reg"] = PROJECTS[pidx].ProjectRegistry
			v1main.Body["regid"] = PROJECTS[pidx].ProjectRegistryId
			v1main.Body["regpw"] = PROJECTS[pidx].ProjectRegistryPw
			v1main.Body["cioption"] = string(yb)

			urecord, err := pkgdbquery.GetUserById(cioption.Process.UserId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while getting user name")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				_ = pkgdbquery.SetProjectCiEndById(newpci.ProjectCiId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			thisKey := urecord.UserName + ":" + electedCluster.ClusterName

			thisAgent, okay := agent_register[thisKey]

			if !okay {
				fail_count += 1

				this_error := fmt.Errorf("error while getting user name")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				_ = pkgdbquery.SetProjectCiEndById(newpci.ProjectCiId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", fmt.Sprintf("failed to find: %s", thisKey), this_error.Error())

				continue

			}

			err = apiximpl.V1ServerPush(v1main, &thisAgent)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while server push")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				_ = pkgdbquery.SetProjectCiEndById(newpci.ProjectCiId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			err = V1PCTL_HandleCiSuccess(&cioption)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while handle ci success")

				_ = V1PCTL_HandleCiError(&cioption, this_error)

				_ = pkgdbquery.SetProjectCiEndById(newpci.ProjectCiId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

		}

		// push cd

		for i := 0; i < CD_Q_LEN; i++ {

			cdoption := CD_Q[i]

			if cdoption.Process.Error {

				err = V1PCTL_HandleCdError(&cdoption, nil)

				if err != nil {

					fail_count += 1

					log.Printf("pctl: skip: error while processing cd error: %s", err.Error())

				}

				continue
			}

			user_clusters, err := pkgdbquery.GetClustersByUserId(cdoption.Process.UserId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while querying user cluster")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			user_projectcds, err := pkgdbquery.GetProjectCdsByProjectId(PROJECTS[cdoption.Process.ProjectIndex].ProjectId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while querying user project cds")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			electedCluster, err := V1PCTL_ElectCdAllocableClusterId(user_clusters, user_projectcds)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while electing cluster for cd")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			var newpcd *pkgresourcedb.DB_Project_CD

			if cdoption.Process.StoredRequest.DependOnCI {

				dependentci, err := pkgdbquery.GetProjectCiRunningByProjectId(PROJECTS[cdoption.Process.ProjectIndex].ProjectId)

				if err != nil {
					fail_count += 1

					this_error := fmt.Errorf("error while getting dependent ci from cd")

					_ = V1PCTL_HandleCdError(&cdoption, this_error)

					log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

					continue

				}

				newpcd, err = pkgdbquery.SetProjectCd(PROJECTS[cdoption.Process.ProjectIndex].ProjectId, electedCluster.ClusterId, dependentci.ProjectCiId)

				if err != nil {
					fail_count += 1

					this_error := fmt.Errorf("error while setting project cd")

					_ = V1PCTL_HandleCdError(&cdoption, this_error)

					log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

					continue

				}

			} else {

				newpcd, err = pkgdbquery.SetProjectCd(PROJECTS[cdoption.Process.ProjectIndex].ProjectId, electedCluster.ClusterId, -1)

				if err != nil {
					fail_count += 1

					this_error := fmt.Errorf("error while setting project cd")

					_ = V1PCTL_HandleCdError(&cdoption, this_error)

					log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

					continue

				}

			}

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while setting project cd")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			lcrecord, err := pkgdbquery.SetLifecycleByProjectId(PROJECTS[cdoption.Process.ProjectIndex].ProjectId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while setting lifecycle")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourcecd.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			cdoption.Process.LifecycleId = lcrecord.LifecycleId

			v1main, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindServerPush, "/project/cluster/cd/alloc", V1MANI)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while getting main")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourcecd.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue
			}

			yb, err := yaml.Marshal(cdoption)

			if err != nil {
				fail_count += 1

				this_error := fmt.Errorf("error while preparing main")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourcecd.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			pidx := cdoption.Process.ProjectIndex

			v1main.Body["name"] = cdoption.Process.ProjectName
			v1main.Body["git"] = PROJECTS[pidx].ProjectGit
			v1main.Body["gitid"] = PROJECTS[pidx].ProjectGitId
			v1main.Body["gitpw"] = PROJECTS[pidx].ProjectGitPw
			v1main.Body["reg"] = PROJECTS[pidx].ProjectRegistry
			v1main.Body["regid"] = PROJECTS[pidx].ProjectRegistryId
			v1main.Body["regpw"] = PROJECTS[pidx].ProjectRegistryPw
			v1main.Body["cdoption"] = string(yb)

			urecord, err := pkgdbquery.GetUserById(cdoption.Process.UserId)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while getting user name")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourceci.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			thisKey := urecord.UserName + ":" + electedCluster.ClusterName

			thisAgent, okay := agent_register[thisKey]

			if !okay {
				fail_count += 1

				this_error := fmt.Errorf("error while getting user name")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourcecd.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", fmt.Sprintf("failed to find: %s", thisKey), this_error.Error())

				continue

			}

			err = apiximpl.V1ServerPush(v1main, &thisAgent)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while server push")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourcecd.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

			err = V1PCTL_HandleCdSuccess(&cdoption)

			if err != nil {

				fail_count += 1

				this_error := fmt.Errorf("error while handle cd success")

				_ = V1PCTL_HandleCdError(&cdoption, this_error)

				_ = pkgdbquery.SetProjectCdEndById(newpcd.ProjectCdId, string(pkgresourcecd.STATUS_ERROR), this_error.Error())

				log.Printf("pctl: %s: %s\n", err.Error(), this_error.Error())

				continue

			}

		}

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

	if SERVER_CONFIG.ResetDBAtStart {

		err := pkgdbquery.DBReset()

		if err != nil {

			return fmt.Errorf("failed to run: reset db: %s", err.Error())
		}

		log.Printf("reset db at start!\n")
	}

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

	listen_addr_client := SERVER_CONFIG.ListenAddr + ":" + SERVER_CONFIG.ListenPortClient

	listen_addr_agent := SERVER_CONFIG.ListenAddr + ":" + SERVER_CONFIG.ListenPortAgent

	go V1ProjectControlLoop()

	go func() {
		http.HandleFunc(SERVER_CONFIG.AgentPath, V1AgentHandler)
		log.Printf("server for agent started at: %s", listen_addr_client+SERVER_CONFIG.AgentPath)
		log.Fatal(http.ListenAndServe(listen_addr_agent, nil))

	}()

	http.HandleFunc(SERVER_CONFIG.ClientPath, V1ClientHandler)

	log.Printf("server for client started at: %s", listen_addr_client+SERVER_CONFIG.ClientPath)

	log.Fatal(http.ListenAndServeTLS(listen_addr_client, "cert-server/cert.pem", "cert-server/cert-priv.pem", nil))

	return nil
}
