package apiximpl

import (
	"fmt"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
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

		if cirecord[i].ProjectCiStatus == pkgresourceci.STATUS_RUNNING {

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

	if status == pkgresourceci.STATUS_RUNNING {

		err = pkgdbquery.SetProjectCiLogById(cirecord[ciidx].ProjectCiId, cilog)

	} else if status == pkgresourceci.STATUS_COMPLETED || status == pkgresourceci.STATUS_ERROR {

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

		if cdrecord[i].ProjectCdStatus == pkgresourceci.STATUS_RUNNING {

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

	if status == pkgresourcecd.STATUS_RUNNING {

		err = pkgdbquery.SetProjectCiLogById(cdrecord[cdidx].ProjectCdId, cdlog)

	} else if status == pkgresourcecd.STATUS_COMPLETED || status == pkgresourcecd.STATUS_ERROR {

		err = pkgdbquery.SetProjectCiEndById(cdrecord[cdidx].ProjectCdId, status, cdlog)

	} else {

		return fmt.Errorf("invalid status: %s", status)
	}

	if err != nil {

		return fmt.Errorf("set: %s", err.Error())
	}

	return nil

}
