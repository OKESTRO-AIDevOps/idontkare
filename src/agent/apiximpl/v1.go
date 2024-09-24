package apiximpl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	pkgcomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/comm"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcecomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/comm"
	pkgresourcelc "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/lifecycle"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

var V1TIMEOUT_MS int = 5000
var interval_ms int = 10

type V1AgentConn struct {
	C          *websocket.Conn
	SessionKey []byte
	Name       string
	UserName   string
	Lock       sync.Mutex
}

type V1AgentCi struct {
	ProjectName string
	Git         string
	GitId       string
	GitPw       string
	Reg         string
	RegId       string
	RegPw       string
	CiOption    pkgresourceci.CiOption
	Status      pkgresourceci.CiStatusType
	Log         string
	BuildLog    string
}

type V1AgentCd struct {
	ProjectName string
	Git         string
	GitId       string
	GitPw       string
	Reg         string
	RegId       string
	RegPw       string
	CdOption    pkgresourcecd.CdOption
	Status      pkgresourcecd.CdStatusType
	Log         string
}

type V1AgentCiCd struct {
	Ci *V1AgentCi `yaml:"ci,omitempty"`
	Cd *V1AgentCd `yaml:"cd,omitempty"`
}

var CI_LOCK sync.Mutex
var CD_LOCK sync.Mutex
var CICD_LOCK sync.Mutex
var LC_LOCK sync.Mutex
var LC_TERM_LOCK sync.Mutex

var CI_OPTIONS_Q = make([]V1AgentCi, 0)
var CD_OPTIONS_Q = make([]V1AgentCd, 0)
var CICD_PIPE_Q = make([]V1AgentCiCd, 0)
var LC_MANIFEST_Q = make([]pkgresourcelc.LifecycleManifest, 0)
var LC_TERMINATOR_Q = make([]pkgresourcelc.LifecycleReport, 0)

func V1CreateCache() error {

	err := os.MkdirAll(".cache", 0777)

	if err != nil {

		return fmt.Errorf("failed to create .cache")
	}

	return nil
}

func V1SaveToCache(filename string, filedata []byte) (string, error) {

	fullPath := ".cache/" + filename

	err := os.WriteFile(fullPath, filedata, 0644)

	if err != nil {

		return "", fmt.Errorf("failed to save to cache: %s", err.Error())
	}

	return fullPath, nil
}

func V1AgentRequestCtl(v1main *pkgresourceapix.V1Main, c *websocket.Conn, sess_key []byte) (*pkgresourceapix.V1ResultData, error) {

	route := v1main.Path

	switch route {

	}

	resp, err := V1RoundTrip(v1main, c, sess_key)

	if err != nil {

		return resp, err
	}

	return resp, nil

}

func V1RoundTrip(v1main *pkgresourceapix.V1Main, c *websocket.Conn, sess_key []byte) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	data, err := yaml.Marshal(v1main)

	if err != nil {

		return nil, fmt.Errorf("roundtrip: marshal: %s", err.Error())
	}

	if sess_key != nil {

		data_str, err := pkgcomm.CommDataEncrypt(data, sess_key)

		if err != nil {

			return nil, fmt.Errorf("roundtrip encrypt: %s", err.Error())
		}

		data = []byte(data_str)

	}
	commjson := pkgresourcecomm.CommJSON{
		Status: pkgresourcecomm.COMM_STATUS_REQUEST,
		Data:   data,
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(interval_ms))

	timelimit_count := 0
	timelimit := V1TIMEOUT_MS / interval_ms

	read_wait_loop := 1

	error_channel := make(chan error)
	read_channel := make(chan *pkgresourcecomm.CommJSON)

	go func() {

		var read_commjson_in pkgresourcecomm.CommJSON

		err := c.ReadJSON(&read_commjson_in)

		if err != nil {

			error_channel <- err
			return
		}

		read_channel <- &read_commjson_in

	}()

	err = c.WriteJSON(commjson)

	if err != nil {
		return nil, fmt.Errorf("routdtrip: write: %s", err.Error())
	}

	var read_commjson *pkgresourcecomm.CommJSON

	for read_wait_loop == 1 {

		select {

		case <-ticker.C:

			timelimit_count += 1

			if timelimit_count >= timelimit {

				ticker.Stop()

				return nil, fmt.Errorf("roundtrip: timeout")

			}

		case read_error := <-error_channel:

			return nil, fmt.Errorf("roundtrip: read: %s", read_error.Error())

		case read_commjson = <-read_channel:

			read_wait_loop = 0

		}
	}

	if read_commjson.Status != pkgresourcecomm.COMM_STATUS_SUCCESS {

		return nil, fmt.Errorf("roundtrip: comm failed: %s", read_commjson.Message)

	}

	resp_data := read_commjson.Data

	if sess_key != nil {

		resp_data, err = pkgcomm.CommDataDecrypt(string(resp_data), sess_key)

		if err != nil {
			return nil, fmt.Errorf("roundtrip: decrypt failed: %s", err.Error())
		}

	}

	err = yaml.Unmarshal(resp_data, &resp)

	if err != nil {

		return nil, fmt.Errorf("roundtrip: data unmarshal: %s", err.Error())
	}

	return &resp, nil
}

func V1ServerPush(v1main *pkgresourceapix.V1Main) error {

	route := v1main.Path

	switch route {

	case "/project/cluster/ci/alloc":

		projectname := v1main.Body["name"]
		git := v1main.Body["git"]
		gitid := v1main.Body["gitid"]
		gitpw := v1main.Body["gitpw"]
		reg := v1main.Body["reg"]
		regid := v1main.Body["regid"]
		regpw := v1main.Body["regpw"]
		cioption := v1main.Body["cioption"]

		err := V1ProjectClusterCiAlloc(projectname, git, gitid, gitpw, reg, regid, regpw, cioption)

		if err != nil {

			return fmt.Errorf("failed server push: ci alloc: %s", err.Error())

		}

	case "/project/cluster/cd/alloc":

		projectname := v1main.Body["name"]
		git := v1main.Body["git"]
		gitid := v1main.Body["gitid"]
		gitpw := v1main.Body["gitpw"]
		reg := v1main.Body["reg"]
		regid := v1main.Body["regid"]
		regpw := v1main.Body["regpw"]
		cdoption := v1main.Body["cdoption"]

		err := V1ProjectClusterCdAlloc(projectname, git, gitid, gitpw, reg, regid, regpw, cdoption)

		if err != nil {

			return fmt.Errorf("failed server push: cd alloc: %s", err.Error())
		}

	case "/lifecycle/manifest/cluster/alloc":

		manifest := v1main.Body["manifest"]

		err := V1LifecycleAlloc(manifest)

		if err != nil {

			return fmt.Errorf("failed server push: lifecycle alloc: %s", err.Error())
		}

	case "/lifecycle/manifest/cluster/free":

		manifest := v1main.Body["manifest"]

		err := V1LifecycleFree(manifest)

		if err != nil {

			return fmt.Errorf("failed server push: lifecycle free: %s", err.Error())
		}

	default:

		return fmt.Errorf("failed server push: no such route: %s", route)

	}

	return nil

}

func V1AgentPush(v1main *pkgresourceapix.V1Main, acon *V1AgentConn) error {

	acon.Lock.Lock()

	defer acon.Lock.Unlock()

	var req pkgresourcecomm.CommJSON

	yb, err := yaml.Marshal(v1main)

	if err != nil {

		return fmt.Errorf("failed agent push: marshal: %s", err.Error())
	}

	enc_b, err := pkgcomm.CommDataEncrypt(yb, acon.SessionKey)

	if err != nil {
		return fmt.Errorf("failed agent push: encrypt: %s", err.Error())
	}

	req.Data = []byte(enc_b)

	err = acon.C.WriteJSON(req)

	if err != nil {

		return fmt.Errorf("failed agent push: write: %s", err.Error())
	}

	return nil

}

func V1CiAdd(ci V1AgentCi) error {

	CI_LOCK.Lock()

	defer CI_LOCK.Unlock()

	CI_OPTIONS_Q = append(CI_OPTIONS_Q, ci)

	return nil
}

func V1CiClear(term_list []int) error {

	CI_LOCK.Lock()

	defer CI_LOCK.Unlock()

	var new_ci_q = make([]V1AgentCi, 0)

	cilen := len(CI_OPTIONS_Q)

	for i := 0; i < cilen; i++ {

		if pkgutils.CheckIfSliceContains[int](term_list, i) {

			continue
		}

		new_ci_q = append(new_ci_q, CI_OPTIONS_Q[i])

	}

	CI_OPTIONS_Q = new_ci_q

	return nil
}

func V1CdAdd(cd V1AgentCd) error {

	CD_LOCK.Lock()

	defer CD_LOCK.Unlock()

	CD_OPTIONS_Q = append(CD_OPTIONS_Q, cd)

	return nil
}

func V1CdClear(term_list []int) error {

	CD_LOCK.Lock()

	defer CD_LOCK.Unlock()

	var new_cd_q = make([]V1AgentCd, 0)

	cdlen := len(CD_OPTIONS_Q)

	for i := 0; i < cdlen; i++ {

		if pkgutils.CheckIfSliceContains[int](term_list, i) {

			continue
		}

		new_cd_q = append(new_cd_q, CD_OPTIONS_Q[i])

	}

	CD_OPTIONS_Q = new_cd_q

	return nil
}

func V1CiCdAdd(ci *V1AgentCi, cd *V1AgentCd) error {

	CICD_LOCK.Lock()

	defer CICD_LOCK.Unlock()

	if ci == nil && cd == nil {

		return fmt.Errorf("update cicd: both empty")
	}

	counterpart_exists := -1

	pipelen := len(CICD_PIPE_Q)

	if ci != nil {

		for i := 0; i < pipelen; i++ {

			if CICD_PIPE_Q[i].Cd != nil {

				if ci.CiOption.Process.ProjectId == CICD_PIPE_Q[i].Cd.CdOption.Process.ProjectId {

					CICD_PIPE_Q[i].Ci = ci

					counterpart_exists = i

					break
				}

			} else {

				continue
			}

		}

		if counterpart_exists == -1 {

			CICD_PIPE_Q = append(CICD_PIPE_Q, V1AgentCiCd{
				Ci: ci,
				Cd: nil,
			})

		}

	} else if cd != nil {

		for i := 0; i < pipelen; i++ {

			if CICD_PIPE_Q[i].Ci != nil {

				if cd.CdOption.Process.ProjectId == CICD_PIPE_Q[i].Ci.CiOption.Process.ProjectId {

					CICD_PIPE_Q[i].Cd = cd

					counterpart_exists = i

					break
				}

			} else {

				continue
			}

		}

		if counterpart_exists == -1 {

			CICD_PIPE_Q = append(CICD_PIPE_Q, V1AgentCiCd{
				Ci: nil,
				Cd: cd,
			})

		}

	}

	return nil
}

func V1CiCdClear(term_list []int) error {
	CICD_LOCK.Lock()

	defer CICD_LOCK.Unlock()

	var new_cicd_q = make([]V1AgentCiCd, 0)

	cicdlen := len(CICD_PIPE_Q)

	for i := 0; i < cicdlen; i++ {

		if pkgutils.CheckIfSliceContains[int](term_list, i) {

			continue
		}

		new_cicd_q = append(new_cicd_q, CICD_PIPE_Q[i])

	}

	CICD_PIPE_Q = new_cicd_q

	return nil
}

func V1LifecycleUpdate(lc_mani *pkgresourcelc.LifecycleManifest) error {

	LC_LOCK.Lock()

	defer LC_LOCK.Unlock()

	var term_idx int = -1

	lclen := len(LC_MANIFEST_Q)

	for i := 0; i < lclen; i++ {

		if lc_mani.Process.ProjectId == LC_MANIFEST_Q[i].Process.ProjectId {

			term_idx = i

			break

		}

	}

	if term_idx != -1 {

		err := V1LifecycleDeleteByIndex(term_idx)

		if err != nil {

			return fmt.Errorf("failed to update: %s", err.Error())

		}
	}

	LC_MANIFEST_Q = append(LC_MANIFEST_Q, *lc_mani)

	return nil
}

func V1LifecycleDelete(lc_mani *pkgresourcelc.LifecycleManifest) error {

	LC_LOCK.Lock()

	LC_TERM_LOCK.Lock()

	defer LC_LOCK.Unlock()

	defer LC_TERM_LOCK.Unlock()

	var term_idx int = -1

	lclen := len(LC_MANIFEST_Q)

	outBuf := bytes.Buffer{}
	errBuf := bytes.Buffer{}

	namespace := lc_mani.Process.ProjectName

	cmd := exec.Command("kubectl", "delete", "namespace", namespace)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err := cmd.Run()

	if err != nil {

		return fmt.Errorf("failed to delete project: %s: %s", namespace, errBuf.String())

	}

	for i := 0; i < lclen; i++ {

		if lc_mani.Process.ProjectId == LC_MANIFEST_Q[i].Process.ProjectId {

			term_idx = i

			break

		}

	}

	if term_idx != -1 {

		err := V1LifecycleDeleteByIndex(term_idx)

		if err != nil {

			return fmt.Errorf("failed to delete: %s", err.Error())

		}
	}

	var report pkgresourcelc.LifecycleReport

	report.Obsolete = true
	report.SentTimestamp = time.Now()
	report.Process = lc_mani.Process

	LC_TERMINATOR_Q = append(LC_TERMINATOR_Q, report)

	return nil
}

func V1LifecycleDeleteByIndex(term_idx int) error {

	var new_lc_q = make([]pkgresourcelc.LifecycleManifest, 0)

	lclen := len(LC_MANIFEST_Q)

	for i := 0; i < lclen; i++ {

		if i == term_idx {

			continue
		}

		new_lc_q = append(new_lc_q, LC_MANIFEST_Q[i])

	}

	LC_MANIFEST_Q = new_lc_q

	return nil
}
