package apiximpl

import (
	"fmt"
	"time"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcelc "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/lifecycle"
	"gopkg.in/yaml.v3"
)

func V1ProjectCiLog(username string, clustername string, projectname string, status string, cilog string) error {

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return fmt.Errorf("get user: %s", err.Error())
	}

	if userrecord == nil {

		return fmt.Errorf("empty user")

	}

	clusterrecord, err := pkgdbquery.GetClustersByUserId(userrecord.UserId)

	if err != nil {

		return fmt.Errorf("cluster: %s", err.Error())
	}

	if clusterrecord == nil {

		return fmt.Errorf("empty cluster")
	}

	clen := len(clusterrecord)
	cidx := -1

	for i := 0; i < clen; i++ {

		if clusterrecord[i].ClusterName == clustername {
			cidx = i
			break
		}

	}

	if cidx == -1 {

		return fmt.Errorf("no such cluster: %s", clustername)
	}

	projectrecord, err := pkgdbquery.GetProjectsByUserId(userrecord.UserId)

	if err != nil {

		return fmt.Errorf("project: %s", err.Error())

	}

	if projectrecord == nil {

		return fmt.Errorf("faield agent push: empty project")
	}

	plen := len(projectrecord)
	pidx := -1

	for i := 0; i < plen; i++ {

		if projectrecord[i].ProjectName == projectname {

			pidx = i

			break
		}

	}

	if pidx == -1 {

		return fmt.Errorf("not found: %s", projectname)
	}

	cirecord, err := pkgdbquery.GetProjectCisByProjectId(projectrecord[pidx].ProjectId)

	if err != nil {

		return fmt.Errorf("project ci: %s", err.Error())

	}

	if cirecord == nil {

		return fmt.Errorf("empty ci record")
	}

	ciidx := -1
	cilen := len(cirecord)

	for i := 0; i < cilen; i++ {

		if !cirecord[i].ProjectCiEnd.Valid {

			ciidx = i

			break
		}

	}

	if ciidx == -1 {

		return fmt.Errorf("running ci not found")
	}

	if cirecord[ciidx].ClusterId != clusterrecord[cidx].ClusterId {

		return fmt.Errorf("cluster id not matching")
	}

	if status == string(pkgresourceci.STATUS_RUNNING) || status == string(pkgresourceci.STATUS_READY) {

		err = pkgdbquery.SetProjectCiLogById(cirecord[ciidx].ProjectCiId, cilog)

	} else if status == string(pkgresourceci.STATUS_COMPLETED) || status == string(pkgresourceci.STATUS_ERROR) {

		err = pkgdbquery.SetProjectCiEndById(cirecord[ciidx].ProjectCiId, status, cilog)

	} else {

		return fmt.Errorf("invalid status: %s", status)
	}

	if err != nil {

		return fmt.Errorf("set: %s", err.Error())
	}

	return nil
}

func V1ProjectCdLog(username string, clustername string, projectname string, status string, cdlog string) error {

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return fmt.Errorf("get user: %s", err.Error())
	}

	if userrecord == nil {

		return fmt.Errorf("empty user")

	}

	clusterrecord, err := pkgdbquery.GetClustersByUserId(userrecord.UserId)

	if err != nil {

		return fmt.Errorf("cluster: %s", err.Error())
	}

	if clusterrecord == nil {

		return fmt.Errorf("empty cluster")
	}

	clen := len(clusterrecord)
	cidx := -1

	for i := 0; i < clen; i++ {

		if clusterrecord[i].ClusterName == clustername {
			cidx = i
			break
		}

	}

	if cidx == -1 {

		return fmt.Errorf("no such cluster: %s", clustername)
	}

	projectrecord, err := pkgdbquery.GetProjectsByUserId(userrecord.UserId)

	if err != nil {

		return fmt.Errorf("project: %s", err.Error())

	}

	if projectrecord == nil {

		return fmt.Errorf("faield agent push: empty project")
	}

	plen := len(projectrecord)
	pidx := -1

	for i := 0; i < plen; i++ {

		if projectrecord[i].ProjectName == projectname {

			pidx = i

			break
		}

	}

	if pidx == -1 {

		return fmt.Errorf("not found: %s", projectname)
	}

	cdrecord, err := pkgdbquery.GetProjectCdsByProjectId(projectrecord[pidx].ProjectId)

	if err != nil {

		return fmt.Errorf("project cd: %s", err.Error())

	}

	if cdrecord == nil {

		return fmt.Errorf("empty cd record")
	}

	cdidx := -1
	cdlen := len(cdrecord)

	for i := 0; i < cdlen; i++ {

		if !cdrecord[i].ProjectCdEnd.Valid {

			cdidx = i

			break
		}

	}

	if cdidx == -1 {

		return fmt.Errorf("running cd not found")
	}

	if cdrecord[cdidx].ClusterId != clusterrecord[cidx].ClusterId {

		return fmt.Errorf("cluster id not matching")
	}

	if status == string(pkgresourcecd.STATUS_RUNNING) || status == string(pkgresourcecd.STATUS_READY) {

		err = pkgdbquery.SetProjectCdLogById(cdrecord[cdidx].ProjectCdId, cdlog)

	} else if status == string(pkgresourcecd.STATUS_COMPLETED) || status == string(pkgresourcecd.STATUS_ERROR) {

		err = pkgdbquery.SetProjectCdEndById(cdrecord[cdidx].ProjectCdId, status, cdlog)

	} else {

		return fmt.Errorf("invalid status: %s", status)
	}

	if err != nil {

		return fmt.Errorf("set: %s", err.Error())
	}

	return nil

}

func V1LifecycleReport(report_str string) error {

	var report pkgresourcelc.LifecycleReport

	report_b := []byte(report_str)

	err := yaml.Unmarshal(report_b, &report)

	if err != nil {

		return fmt.Errorf("failed to unmarshal report: %s", err.Error())
	}

	report.ReceivedTimestamp = time.Now()

	report_b, err = yaml.Marshal(report)

	if err != nil {

		return fmt.Errorf("failed to add recvd time: %s", err.Error())
	}

	err = pkgdbquery.SetLifecycleReportByLifecycleId(report.Process.LifecylcleId, string(report_b))

	if err != nil {

		return fmt.Errorf("failed to set lc report: %s", err.Error())
	}

	return nil
}
