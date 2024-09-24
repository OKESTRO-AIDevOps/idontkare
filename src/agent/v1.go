package main

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	"github.com/OKESTRO-AIDevOps/idontkare/pkg/comm"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourceauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/auth"
	pkgresourcecomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/comm"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	agentapiximpl "github.com/OKESTRO-AIDevOps/idontkare/src/agent/apiximpl"
	"github.com/goccy/go-yaml"
	"github.com/gorilla/websocket"
)

func V1SetApixImpl() {

	agentapiximpl.V1TIMEOUT_MS = AGENT_CONFIG.TimeoutMS
}

func V1Connect(connect_url string, name string, username string, priv_key *rsa.PrivateKey, mani *pkgresourceapix.V1Manifest) (*agentapiximpl.V1AgentConn, error) {

	var agentConn agentapiximpl.V1AgentConn

	v1mainTmpl, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindAgentRequestPriv, "/cluster/connect", mani)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: copy main: %s", err.Error())
	}

	if v1mainTmpl == nil {

		return nil, fmt.Errorf("failed to connect: empty copy: %s", err.Error())
	}

	v1mainTmpl.Body["name"] = name
	v1mainTmpl.Body["username"] = username

	yb, err := yaml.Marshal(v1mainTmpl)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: marshal: %s", err.Error())
	}

	v1main, err := pkgapix.V1GetMainByByte(yb, mani)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: get valid main: %s", err.Error())
	}

	c, _, err := websocket.DefaultDialer.Dial(connect_url, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: dial: %s", err.Error())
	}

	resp, err := agentapiximpl.V1AgentRequestCtl(v1main, c, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: get chal: %s", err.Error())
	}

	chal_key_enc_b, err := hex.DecodeString(resp.Output)

	if err != nil {
		return nil, fmt.Errorf("failed to connect: decode: %s", err.Error())
	}

	chal_key_b, err := pkgutils.DecryptWithPrivateKey(chal_key_enc_b, priv_key)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: decrypt: %s", err.Error())
	}

	session_key, _ := pkgutils.RandomHex(16)

	chal_data := pkgresourceauth.ChallengeData{
		Pass: AGENT_CONFIG.UserPass,
		Key:  session_key,
	}

	chal_data_b, err := json.Marshal(chal_data)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: marshal: %s", err.Error())
	}

	session_key_b := []byte(session_key)

	data_b, err := pkgutils.EncryptWithSymmetricKey(chal_key_b, chal_data_b)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: encrypt chal: %s", err.Error())
	}

	data_hex := hex.EncodeToString(data_b)

	v1mainChalTmpl, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindAgentRequestPriv, "/cluster/connect/challenge", mani)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: copy main for chal: %s", err.Error())
	}

	v1mainChalTmpl.Body["chaldata"] = data_hex

	yb, err = yaml.Marshal(v1mainChalTmpl)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: main challenge marshal: %s", err.Error())
	}

	v1mainChal, err := pkgapix.V1GetMainByByte(yb, mani)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: get main challenge: %s", err.Error())
	}

	// TODO:
	//   should check resp
	//   decrypt with session_key_b to check
	//   expected response output message

	_, err = agentapiximpl.V1AgentRequestCtl(v1mainChal, c, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: send main challenge: %s", err.Error())
	}

	agentConn.C = c
	agentConn.Name = name
	agentConn.UserName = username
	agentConn.SessionKey = session_key_b

	return &agentConn, nil
}

func V1HandleAgentPush(acon *agentapiximpl.V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	go agentapiximpl.V1CiHandler(acon, mani)

	log.Printf("started ci handler\n")

	go agentapiximpl.V1CdHandler(acon, mani)

	log.Printf("started cd handler\n")

	go agentapiximpl.V1CiCdHandler(acon, mani)

	log.Printf("started cicd handler\n")

	go agentapiximpl.V1LifecycleHandler(acon, mani)

	log.Printf("started lifecycle handler\n")

	go agentapiximpl.V1LifecycleTerminator(acon, mani)

	log.Printf("started lifecycle terminator\n")
}

func V1HandleServerPush(acon *agentapiximpl.V1AgentConn, mani *pkgresourceapix.V1Manifest) error {

	var errMessage error

	for {

		req := pkgresourcecomm.CommJSON{}

		err := acon.C.ReadJSON(&req)

		if err != nil {

			errMessage = fmt.Errorf("failed to communicate: read: %s", err.Error())

			log.Println(errMessage)

			break
		}

		data_b, err := comm.CommDataDecrypt(string(req.Data), acon.SessionKey)

		if err != nil {

			errMessage = fmt.Errorf("failed to communicate: decrypt: %s", err.Error())

			log.Println(errMessage)

			continue
		}

		v1main, err := pkgapix.V1GetMainByByte(data_b, mani)

		if err != nil {

			errMessage = fmt.Errorf("failed to communicate: get main: %s", err.Error())

			log.Println(errMessage)

			continue
		}

		log.Println(v1main.Kind + v1main.Path)

		switch v1main.Kind {

		case pkgresourceapix.V1KindServerPush:

			err = agentapiximpl.V1ServerPush(v1main)

			if err != nil {

				errMessage = fmt.Errorf("failed to handle write: %s", err.Error())

				log.Println(errMessage)

				continue
			}

		default:

			errMessage = fmt.Errorf("failed to communicate: illegal kind: %s", v1main.Kind)

			log.Println(errMessage)

			continue
		}

		//log.Printf("handled server push: %s: %s", v1main.Kind, v1main.Path)

	}

	return nil
}

func V1Run(connect_url string, name string, username string, key_path string) error {

	V1SetApixImpl()

	err := agentapiximpl.V1CreateCache()

	if err != nil {
		return fmt.Errorf("failed to run: cache: %s", err.Error())
	}

	manifest, err := pkgapix.V1GetManifest()

	if err != nil {

		return fmt.Errorf("failed to run: manifest: %s", err.Error())
	}

	key_file_b, err := os.ReadFile(key_path)

	if err != nil {

		return fmt.Errorf("failed to run: read key: %s", err.Error())
	}

	privKey, err := pkgutils.BytesToPrivateKey(key_file_b)

	if err != nil {
		return fmt.Errorf("failed to run: parse key: %s", err.Error())
	}

	conn, err := V1Connect(connect_url, name, username, privKey, manifest)

	if err != nil {

		return err
	}

	V1HandleAgentPush(conn, manifest)

	err = V1HandleServerPush(conn, manifest)

	if err != nil {

		return err
	}

	return nil
}
