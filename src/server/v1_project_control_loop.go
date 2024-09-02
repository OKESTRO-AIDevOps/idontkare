package main

import (
	"fmt"
	"time"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	"gopkg.in/yaml.v3"
)

func V1PCTL_HandleCiSuccess(cioption *pkgresourceci.CiOption) error {

	opt_byte := []byte{}

	cioption.Response = &struct {
		ProcessedTimestamp time.Time "yaml:\"processed_timestamp\""
		Error              bool      "yaml:\"error\""
		Log                string    `yaml:"log"`
	}{
		ProcessedTimestamp: time.Now(),
		Error:              false,
		Log:                "success",
	}

	yb, err := yaml.Marshal(cioption)

	if err != nil {

		yb = opt_byte
	}

	err = pkgdbquery.SetProjectCiOptionByUserIdAndName(cioption.Process.UserId, cioption.Process.ProjectName, string(yb))

	if err != nil {

		return fmt.Errorf("faield to handle ci success: %s", err.Error())
	}

	return nil
}

func V1PCTL_HandleCiError(cioption *pkgresourceci.CiOption, resp_err error) error {

	opt_byte := []byte{}

	if resp_err == nil {
		cioption.Response = &struct {
			ProcessedTimestamp time.Time "yaml:\"processed_timestamp\""
			Error              bool      "yaml:\"error\""
			Log                string    `yaml:"log"`
		}{
			ProcessedTimestamp: time.Now(),
			Error:              cioption.Process.Error,
			Log:                cioption.Process.Log,
		}

	} else {
		cioption.Response = &struct {
			ProcessedTimestamp time.Time "yaml:\"processed_timestamp\""
			Error              bool      "yaml:\"error\""
			Log                string    `yaml:"log"`
		}{
			ProcessedTimestamp: time.Now(),
			Error:              true,
			Log:                resp_err.Error(),
		}

	}

	yb, err := yaml.Marshal(cioption)

	if err != nil {

		yb = opt_byte
	}

	err = pkgdbquery.SetProjectCiOptionByUserIdAndName(cioption.Process.UserId, cioption.Process.ProjectName, string(yb))

	if err != nil {

		return fmt.Errorf("faield to handle ci error: %s", err.Error())
	}

	return nil
}

func V1PCTL_HandleCdSuccess(cdoption *pkgresourcecd.CdOption) error {

	opt_byte := []byte{}

	cdoption.Response = &struct {
		ProcessedTimestamp time.Time "yaml:\"processed_timestamp\""
		Error              bool      "yaml:\"error\""
		Log                string    `yaml:"log"`
	}{
		ProcessedTimestamp: time.Now(),
		Error:              false,
		Log:                "success",
	}

	yb, err := yaml.Marshal(cdoption)

	if err != nil {

		yb = opt_byte
	}

	err = pkgdbquery.SetProjectCdOptionByUserIdAndName(cdoption.Process.UserId, cdoption.Process.ProjectName, string(yb))

	if err != nil {

		return fmt.Errorf("faield to handle cd success: %s", err.Error())
	}

	return nil
}

func V1PCTL_HandleCdError(cdoption *pkgresourcecd.CdOption, resp_err error) error {

	opt_byte := []byte{}

	if resp_err == nil {

		cdoption.Response = &struct {
			ProcessedTimestamp time.Time "yaml:\"processed_timestamp\""
			Error              bool      "yaml:\"error\""
			Log                string    `yaml:"log"`
		}{
			ProcessedTimestamp: time.Now(),
			Error:              cdoption.Process.Error,
			Log:                cdoption.Process.Log,
		}

	} else {
		cdoption.Response = &struct {
			ProcessedTimestamp time.Time "yaml:\"processed_timestamp\""
			Error              bool      "yaml:\"error\""
			Log                string    `yaml:"log"`
		}{
			ProcessedTimestamp: time.Now(),
			Error:              true,
			Log:                resp_err.Error(),
		}

	}

	yb, err := yaml.Marshal(cdoption)

	if err != nil {

		yb = opt_byte
	}

	err = pkgdbquery.SetProjectCdOptionByUserIdAndName(cdoption.Process.UserId, cdoption.Process.ProjectName, string(yb))

	if err != nil {

		return fmt.Errorf("failed to handle cd error: %s", err.Error())
	}

	return nil
}

func V1PCTL_ElectCiAllocableClusterId(user_clusters []pkgresourcedb.DB_Cluster, user_projectcis []pkgresourcedb.DB_Project_CI) (*pkgresourcedb.DB_Cluster, error) {

	var electedCluster *pkgresourcedb.DB_Cluster

	var validCluster []pkgresourcedb.DB_Cluster

	uclen := len(user_clusters)

	if uclen < 1 {

		return nil, fmt.Errorf("failed to elect: no user cluster")
	}

	for i := 0; i < uclen; i++ {

		if user_clusters[i].ClusterConnected == 1 {

			validCluster = append(validCluster, user_clusters[i])
		}

	}

	vclen := len(validCluster)

	if vclen == 0 {

		return nil, fmt.Errorf("failed to elect: no valid cluster exists")
	}

	uplen := len(user_projectcis)

	for i := 0; i < uplen; i++ {

		if !user_projectcis[i].ProjectCiEnd.Valid {

			return nil, fmt.Errorf("failed to elect: ci already in process")

		}

	}

	picked := pkgutils.GetRandIntInRange(0, uclen-1)

	electedCluster = &validCluster[picked]

	return electedCluster, nil
}

func V1PCTL_ElectCdAllocableClusterId(user_clusters []pkgresourcedb.DB_Cluster, user_projectcds []pkgresourcedb.DB_Project_CD) (*pkgresourcedb.DB_Cluster, error) {

	var electedCluster *pkgresourcedb.DB_Cluster

	var validCluster []pkgresourcedb.DB_Cluster

	uclen := len(user_clusters)

	if uclen < 1 {

		return nil, fmt.Errorf("failed to elect: no user cluster")
	}

	for i := 0; i < uclen; i++ {

		if user_clusters[i].ClusterConnected == 1 {

			validCluster = append(validCluster, user_clusters[i])
		}

	}

	vclen := len(validCluster)

	if vclen == 0 {

		return nil, fmt.Errorf("failed to elect: no valid cluster exists")
	}

	uplen := len(user_projectcds)

	for i := 0; i < uplen; i++ {

		if !user_projectcds[i].ProjectCdEnd.Valid {

			return nil, fmt.Errorf("failed to elect: cd already in process")

		}

	}

	picked := pkgutils.GetRandIntInRange(0, uclen-1)

	electedCluster = &validCluster[picked]

	return electedCluster, nil

}
