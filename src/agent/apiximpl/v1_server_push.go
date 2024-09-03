package apiximpl

import (
	"fmt"

	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	"gopkg.in/yaml.v3"
)

func V1ProjectClusterCiAlloc(projectname string, git string, gitid string, gitpw string, reg string, regid string, regpw string, cioption string) error {

	// TODO:
	//   start building

	var agentCi V1AgentCi
	var ciOption pkgresourceci.CiOption

	cioption_b := []byte(cioption)

	err := yaml.Unmarshal(cioption_b, &ciOption)

	if err != nil {

		return fmt.Errorf("failed to unmarshal ci option: %s", err.Error())
	}

	agentCi.ProjectName = projectname
	agentCi.Git = git
	agentCi.GitId = gitid
	agentCi.GitPw = gitpw
	agentCi.Reg = reg
	agentCi.RegId = regid
	agentCi.RegPw = regpw
	agentCi.CiOption = ciOption
	agentCi.Status = pkgresourceci.STATUS_READY

	if agentCi.CiOption.Process.LinkToCd != -1 {

		err := V1CiCdAdd(&agentCi, nil)

		if err != nil {

			return fmt.Errorf("failed to add to cicd as link exists: %s", err.Error())
		}

	} else {

		err := V1CiAdd(agentCi)

		if err != nil {

			return fmt.Errorf("faield to add to ci as link doesn't exist: %s", err.Error())
		}
	}

	return nil
}

func V1ProjectClusterCdAlloc(projectname string, git string, gitid string, gitpw string, reg string, regid string, regpw string, cdoption string) error {

	// TODO:
	//   start deploying

	var agentCd V1AgentCd
	var cdOption pkgresourcecd.CdOption

	cdoption_b := []byte(cdoption)

	err := yaml.Unmarshal(cdoption_b, &cdOption)

	if err != nil {

		return fmt.Errorf("failed to unmarshal cd option: %s", err.Error())
	}

	agentCd.ProjectName = projectname
	agentCd.Git = git
	agentCd.GitId = gitid
	agentCd.GitPw = gitpw
	agentCd.Reg = reg
	agentCd.RegId = regid
	agentCd.RegPw = regpw
	agentCd.CdOption = cdOption
	agentCd.Status = pkgresourcecd.STATUS_READY

	if agentCd.CdOption.Process.StoredRequest.DependOnCI {

		err := V1CiCdAdd(nil, &agentCd)

		if err != nil {

			return fmt.Errorf("failed to add cicd dependent on ci: %s", err.Error())
		}

	} else {

		err := V1CdAdd(agentCd)

		if err != nil {

			return fmt.Errorf("failed to add cd independent of ci: %s", err.Error())
		}

	}

	return nil
}
