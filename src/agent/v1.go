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

type V1AgentConn struct {
	C          *websocket.Conn
	SessionKey []byte
	Name       string
	UserName   string
}

func V1SetApixImpl() {

	agentapiximpl.V1TIMEOUT_MS = AGENT_CONFIG.TimeoutMS
}

func V1Connect(connect_url string, name string, username string, priv_key *rsa.PrivateKey, mani *pkgresourceapix.V1Manifest) (*V1AgentConn, error) {

	var agentConn V1AgentConn

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

	session_key, _ := pkgutils.RandomHex(32)

	chal_data := pkgresourceauth.ChallengeData{
		Pass: AGENT_CONFIG.UserPass,
		Key:  session_key,
	}

	chal_data_b, err := json.Marshal(chal_data)

	if err != nil {

		return nil, fmt.Errorf("failed to connect: marshal: %s", err.Error())
	}

	session_key_b, _ := hex.DecodeString(session_key)

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

func V1Communicate(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) error {

	keep_comm := 1

	for keep_comm == 1 {

		var retData []byte

		var errMessage error = nil

		req := pkgresourcecomm.CommJSON{}

		resp := pkgresourcecomm.CommJSON{}

		for {

			err := acon.C.ReadJSON(&req)

			if err != nil {

				errMessage = fmt.Errorf("failed to communicate: read: %s", err.Error())

				break
			}

			data_b, err := comm.CommDataDecrypt(string(req.Data), acon.SessionKey)

			if err != nil {

				errMessage = fmt.Errorf("failed to communicate: decrypt: %s", err.Error())

				break
			}

			v1main, err := pkgapix.V1GetMainByByte(data_b, mani)

			if err != nil {

				errMessage = fmt.Errorf("failed to communicate: get main: %s", err.Error())

				break
			}

			var respData *pkgresourceapix.V1ResultData

			switch v1main.Kind {

			case pkgresourceapix.V1KindServerWrite:

				respData, err = agentapiximpl.V1ServerWriteCtl(v1main)

				if err != nil {

					errMessage = fmt.Errorf("failed to handle write: %s", err.Error())

					break
				}

			case pkgresourceapix.V1KindServerRead:

				respData, err = agentapiximpl.V1ServerReadCtl(v1main)

				if err != nil {

					errMessage = fmt.Errorf("failed to handle read: %s", err.Error())

					break
				}

			default:

				errMessage = fmt.Errorf("failed to communicate: illegal kind: %s", v1main.Kind)

				break
			}

			if errMessage != nil {
				break
			}

			rd_b, err := yaml.Marshal(*respData)

			if err != nil {

				errMessage = fmt.Errorf("failed to communicate: marshal data: %s", err.Error())

				break
			}

			data_enc, err := comm.CommDataEncrypt(rd_b, acon.SessionKey)

			if err != nil {

				errMessage = fmt.Errorf("failed to communicate: encrypt: %s", err.Error())

				break
			}

			retData = []byte(data_enc)

			break
		}

		if errMessage != nil {

			log.Printf("error communicate: %s\n", errMessage.Error())

			resp.Data = []byte{}
			resp.Status = pkgresourcecomm.COMM_STATUS_FAILURE

		} else {

			resp.Data = retData
			resp.Status = pkgresourcecomm.COMM_STATUS_SUCCESS

		}

		err := acon.C.WriteJSON(resp)

		if err != nil {

			return fmt.Errorf("failed to write: %s", err.Error())
		}

	}

	return nil
}

func V1Run(connect_url string, name string, username string, key_path string) error {

	V1SetApixImpl()

	manifest, err := pkgapix.V1GetManifest()

	if err != nil {

		return err
	}

	key_file_b, err := os.ReadFile(key_path)

	if err != nil {

		return err
	}

	privKey, err := pkgutils.BytesToPrivateKey(key_file_b)

	if err != nil {
		return err
	}

	conn, err := V1Connect(connect_url, name, username, privKey, manifest)

	if err != nil {

		return err
	}

	err = V1Communicate(conn, manifest)

	if err != nil {

		return err
	}

	return nil
}
