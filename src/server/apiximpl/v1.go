package apiximpl

import (
	"fmt"
	"log"

	"github.com/OKESTRO-AIDevOps/idontkare/pkg/comm"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourceauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/auth"
	pkgcomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/comm"
	"gopkg.in/yaml.v3"
)

func V1ClientRequestCtl(v1main *pkgresourceapix.V1Main) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	route := v1main.Path

	switch route {

	case "/user/set":

		name := v1main.Body["name"]
		pass := v1main.Body["pass"]

		err := V1UserSet(name, pass)

		if err != nil {

			log.Printf("failed to set user: %s: %s", name, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE
			resp.Output = fmt.Sprintf("failed to set user: %s", name)

		} else {
			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Output = fmt.Sprintf("successfully set user: %s", name)

		}

	case "/cluster/set":

		username := v1main.Body["username"]
		name := v1main.Body["name"]

		privstr, err := V1ClusterSet(username, name)

		if err != nil {

			log.Printf("failed to set cluster: %s: %s", name, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to set cluster: %s", name)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Output = privstr
		}

	case "/project/set":

		username := v1main.Body["username"]
		projectname := v1main.Body["name"]
		git := v1main.Body["git"]
		gitid := v1main.Body["gitid"]
		gitpw := v1main.Body["gitpw"]
		reg := v1main.Body["reg"]
		regid := v1main.Body["regid"]
		regpw := v1main.Body["regpw"]

		err := V1ProjectSet(username, projectname, git, gitid, gitpw, reg, regid, regpw)

		if err != nil {

			log.Printf("failed to set project: %s: %s", projectname, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to set project: %s", projectname)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Output = fmt.Sprintf("successfully set project: %s", projectname)
		}

	case "/project/ci/option/set":

		username := v1main.Body["username"]
		projectname := v1main.Body["name"]
		cioptiondata := v1main.Body["path"]

		err := V1ProjectCiOptionSet(username, projectname, cioptiondata)

		if err != nil {

			log.Printf("failed to set project ci option: %s: %s", projectname, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to set project ci option: %s", projectname)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Output = fmt.Sprintf("successfully set project ci option: %s", projectname)
		}

	case "/project/cd/option/set":

		username := v1main.Body["username"]
		projectname := v1main.Body["name"]
		cdoptiondata := v1main.Body["path"]

		err := V1ProjectCdOptionSet(username, projectname, cdoptiondata)

		if err != nil {

			log.Printf("failed to set project cd option: %s: %s", projectname, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to set project cd option: %s", projectname)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Output = fmt.Sprintf("successfully set project cd option: %s", projectname)
		}

	case "/project/ci/history/get/all":

		username := v1main.Body["username"]
		projectname := v1main.Body["project"]

		out, err := V1ProjectCiHistoryGetAll(username, projectname)

		if err != nil {

			log.Printf("failed to get all project ci history: %s: %s", projectname, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to get all project ci history: %s", projectname)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Format = pkgresourceapix.V1RESULT_FORMAT_YAML
			resp.Output = out
		}

	case "/project/cd/history/get/all":

		username := v1main.Body["username"]
		projectname := v1main.Body["project"]

		out, err := V1ProjectCdHistoryGetAll(username, projectname)

		if err != nil {

			log.Printf("failed to get all project cd history: %s: %s", projectname, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to get all project cd history: %s", projectname)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Format = pkgresourceapix.V1RESULT_FORMAT_YAML
			resp.Output = out
		}

	case "/lifecycle/report/get/latest":

		username := v1main.Body["username"]
		projectname := v1main.Body["project"]

		out, err := V1LifecycleReportGetLatest(username, projectname)

		if err != nil {

			log.Printf("failed to get latest lifecycle report: %s: %s", projectname, err.Error())

			resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE

			resp.Output = fmt.Sprintf("failed to get latest lifecycle report: %s", projectname)

		} else {

			resp.Status = pkgresourceapix.V1RESULT_STATUS_SUCCESS
			resp.Format = pkgresourceapix.V1RESULT_FORMAT_YAML
			resp.Output = out
		}

	default:

		resp.Status = pkgresourceapix.V1RESULT_STATUS_FAILURE
		resp.Output = "no such path: " + route

	}

	return &resp, nil

}

func V1AgentPush(v1main *pkgresourceapix.V1Main, username string, clustername string) error {

	route := v1main.Path

	switch route {

	case "/project/ci/log":

		projectname := v1main.Body["name"]
		status := v1main.Body["status"]
		cilog := v1main.Body["log"]

		err := V1ProjectCiLog(username, clustername, projectname, status, cilog)

		if err != nil {

			return fmt.Errorf("agent push: project ci log: %s", err.Error())
		}

	case "/project/cd/log":

		projectname := v1main.Body["name"]
		status := v1main.Body["status"]
		cdlog := v1main.Body["log"]

		err := V1ProjectCdLog(username, clustername, projectname, status, cdlog)

		if err != nil {

			return fmt.Errorf("agent push: project cd log: %s", err.Error())
		}

	case "/project/lifecycle/report":

		lc_report := v1main.Body["report"]

		err := V1LifecycleReport(lc_report)

		if err != nil {

			return fmt.Errorf("agent push: lifecycle report: %s", err.Error())
		}

	default:

		return fmt.Errorf("failed agent push: no such route: %s", route)
	}

	return nil

}

func V1ServerPush(v1main *pkgresourceapix.V1Main, agent *pkgresourceauth.AgentData) error {

	var req pkgcomm.CommJSON

	yb, err := yaml.Marshal(v1main)

	if err != nil {

		return fmt.Errorf("failed server push: marshal: %s", err.Error())
	}

	log.Println(v1main.Kind + v1main.Path)

	enc_b, err := comm.CommDataEncrypt(yb, []byte(agent.Key))

	if err != nil {

		return fmt.Errorf("failed server push: encrypt: %s", err.Error())
	}

	req.Data = []byte(enc_b)

	err = agent.C.WriteJSON(req)

	if err != nil {

		return fmt.Errorf("failed server push: write: %s", err.Error())
	}

	return nil
}
