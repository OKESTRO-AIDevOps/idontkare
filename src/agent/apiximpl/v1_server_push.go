package apiximpl

import (
	"fmt"

	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcelc "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/lifecycle"
	"gopkg.in/yaml.v3"
)

func V1ProjectClusterCiAlloc(projectname string, git string, gitid string, gitpw string, reg string, regid string, regpw string, cioption string) error {

	// TODO:
	//   start building

	var ciOption pkgresourceci.CiOption

	cioption_b := []byte(cioption)

	err := yaml.Unmarshal(cioption_b, &ciOption)

	if err != nil {

		return fmt.Errorf("failed to unmarshal ci option: %s", err.Error())
	}

	CI_OPTIONS_Q = append(CI_OPTIONS_Q, ciOption)

	return nil
}

func V1ProjectClusterCdAlloc(projectname string, git string, gitid string, gitpw string, reg string, regid string, regpw string, cdoption string) error {

	// TODO:
	//   start deploying

	var cdOption pkgresourcecd.CdOption

	cdoption_b := []byte(cdoption)

	err := yaml.Unmarshal(cdoption_b, &cdOption)

	if err != nil {

		return fmt.Errorf("failed to unmarshal cd option: %s", err.Error())
	}

	CD_OPTIONS_Q = append(CD_OPTIONS_Q, cdOption)

	return nil
}

func V1ProjectLifecycleUpdate(projectname string, lcoption string) error {

	// TODO:
	//  check if lifecycle exists

	var lcOption pkgresourcelc.LifecycleOption

	lcoption_b := []byte(lcoption)

	err := yaml.Unmarshal(lcoption_b, &lcOption)

	if err != nil {

		return fmt.Errorf("failed to unmarshal lc option: %s", err.Error())
	}

	LC_OPTIONS_Q = append(LC_OPTIONS_Q, lcOption)

	return nil
}
