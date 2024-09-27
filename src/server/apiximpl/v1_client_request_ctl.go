package apiximpl

import (
	"fmt"
	"strings"
	"time"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	"gopkg.in/yaml.v3"
)

func V1UserSet(name string, pass string) error {

	err := pkgdbquery.SetUser(name, pass)

	if err != nil {

		return fmt.Errorf("failed to set user: %s", err.Error())
	}

	return nil

}

func V1ClusterSet(username string, name string) (string, error) {

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return "", fmt.Errorf("failed to set cluster: %s", err.Error())
	}

	if userrecord == nil {

		return "", fmt.Errorf("failed to set cluster: empty record")
	}

	priv, pub, err := pkgutils.GenerateKeyPair(4096)

	if err != nil {

		return "", fmt.Errorf("failed to set cluster: generate: %s", err.Error())
	}

	priv_b, err := pkgutils.PrivateKeyToBytes(priv)

	if err != nil {

		return "", fmt.Errorf("failed to set cluster: priv b: %s", err.Error())
	}

	pub_b, err := pkgutils.PublicKeyToBytes(pub)

	if err != nil {

		return "", fmt.Errorf("failed to set cluster: pub b: %s", err.Error())
	}

	priv_pem := string(priv_b)

	pub_pem := string(pub_b)

	err = pkgdbquery.SetCluster(userrecord.UserId, name, pub_pem)

	if err != nil {

		return "", fmt.Errorf("failed to set cluster: set cluster failed: %s", err.Error())
	}

	return priv_pem, nil

}

func V1ProjectSet(username string, projectname string, git string, gitid string, gitpw string, reg string, regid string, regpw string) error {

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return fmt.Errorf("failed to set project: %s", err.Error())
	}

	if userrecord == nil {

		return fmt.Errorf("failed to set project: empty record")
	}

	checkproto := strings.HasPrefix(git, "https://")

	if checkproto {

		gitnew := strings.ReplaceAll(git, "https://", "")

		git = gitnew
	}

	checkproto = strings.HasPrefix(git, "http://")

	if checkproto {

		return fmt.Errorf("failed to set project: git http:// not allowed")

	}

	checkproto = strings.HasPrefix(reg, "https://")

	if checkproto {

		regnew := strings.ReplaceAll(reg, "https://", "")

		reg = regnew
	}

	checkproto = strings.HasPrefix(reg, "http://")

	if checkproto {

		return fmt.Errorf("failed to set project: reg http:// not allowed")
	}

	err = pkgdbquery.SetProject(
		userrecord.UserId,
		projectname,
		git,
		gitid,
		gitpw,
		reg,
		regid,
		regpw,
	)

	if err != nil {

		return fmt.Errorf("failed to set project: set: %s", err.Error())
	}

	return nil
}

func V1ProjectCiOptionSet(username string, projectname string, cioptiondata string) error {

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return fmt.Errorf("failed to set project ci option: %s", err.Error())
	}

	if userrecord == nil {

		return fmt.Errorf("failed to set project ci option: empty record")
	}

	err = pkgdbquery.SetProjectCiOptionByUserIdAndName(userrecord.UserId, projectname, cioptiondata)

	if err != nil {

		return fmt.Errorf("failed to set project ci option: set: %s", err.Error())
	}

	return nil
}

func V1ProjectCdOptionSet(username string, projectname string, cdoptiondata string) error {

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return fmt.Errorf("failed to set project cd option: %s", err.Error())
	}

	if userrecord == nil {

		return fmt.Errorf("failed to set project cd option: empty record")
	}

	err = pkgdbquery.SetProjectCdOptionByUserIdAndName(userrecord.UserId, projectname, cdoptiondata)

	if err != nil {

		return fmt.Errorf("failed to set project cd option: set: %s", err.Error())
	}

	return nil
}

func V1ProjectCiHistoryGetAll(username string, projectname string) (string, error) {

	var output string

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return "", fmt.Errorf("failed to get ci histall: no such user: %s", err.Error())
	}

	if userrecord == nil {

		return "", fmt.Errorf("failed to get ci histall: empty user record")
	}

	projectrecords, err := pkgdbquery.GetProjectsByUserId(userrecord.UserId)

	if err != nil {

		return "", fmt.Errorf("failed to get ci histall: failed to get projects for: %s", username)
	}

	plen := len(projectrecords)

	if plen == 0 {

		return "", fmt.Errorf("failed to get ci histall: no associated projects for: %s", username)
	}

	idx := -1

	for i := 0; i < plen; i++ {

		if projectrecords[i].ProjectName == projectname {

			idx = i

			break

		}
	}

	if idx == -1 {

		return "", fmt.Errorf("failed to get ci histall: no such project: %s", projectname)
	}

	cis, err := pkgdbquery.GetProjectCisByProjectId(projectrecords[idx].ProjectId)

	if err != nil {

		return "", fmt.Errorf("failed to get cis: %s", err.Error())

	}

	cilen := len(cis)

	if cilen == 0 {

		return "", fmt.Errorf("failed to get cis: len: %d", cilen)
	}

	clusterrecords, err := pkgdbquery.GetClustersByUserId(userrecord.UserId)

	if err != nil {

		return "", fmt.Errorf("failed to get clusters: %s", err.Error())
	}

	clen := len(clusterrecords)

	if clen == 0 {

		return "", fmt.Errorf("failed to get clusters: %d", clen)
	}

	out_struct := make([]pkgresourceapix.V1ProjectCiExport, 0)

	for i := 0; i < cilen; i++ {

		tmpout := pkgresourceapix.V1ProjectCiExport{}

		tmpout.ProjectName = projectname

		cfound := -1

		for j := 0; j < clen; j++ {

			if cis[i].ClusterId == clusterrecords[j].ClusterId {

				cfound = j

				break

			}

		}

		if cfound == -1 {

			tmpout.ClusterName = pkgresourceapix.V1EXPORT_STRING_UNKNOWN
		} else {

			tmpout.ClusterName = clusterrecords[cfound].ClusterName
		}

		tmpout.ProjectCiStatus = cis[i].ProjectCiStatus

		if !cis[i].ProjectCiLog.Valid {

			tmpout.ProjectCiLog = pkgresourceapix.V1EXPORT_STRING_NULL

		} else {

			tmpout.ProjectCiLog = cis[i].ProjectCiLog.String

		}

		tmpout.ProjectCiStart = cis[i].ProjectCiStart

		if !cis[i].ProjectCiEnd.Valid {

			tmpout.ProjectCiEnd = time.Time{}

		} else {
			tmpout.ProjectCiEnd = cis[i].ProjectCiEnd.Time
		}

		out_struct = append(out_struct, tmpout)

	}

	outb, err := yaml.Marshal(out_struct)

	if err != nil {

		return "", fmt.Errorf("failed to marshal ci history: %s", err.Error())
	}

	output = string(outb)

	return output, nil
}

func V1ProjectCdHistoryGetAll(username string, projectname string) (string, error) {

	var output string

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return "", fmt.Errorf("failed to get cd histall: no such user: %s", err.Error())
	}

	if userrecord == nil {

		return "", fmt.Errorf("failed to get cd histall: empty user record")
	}

	projectrecords, err := pkgdbquery.GetProjectsByUserId(userrecord.UserId)

	if err != nil {

		return "", fmt.Errorf("failed to get cd histall: failed to get projects for: %s", username)
	}

	plen := len(projectrecords)

	if plen == 0 {

		return "", fmt.Errorf("failed to get cd histall: no associated projects for: %s", username)
	}

	idx := -1

	for i := 0; i < plen; i++ {

		if projectrecords[i].ProjectName == projectname {

			idx = i

			break

		}
	}

	if idx == -1 {

		return "", fmt.Errorf("failed to get cd histall: no such project: %s", projectname)
	}

	cds, err := pkgdbquery.GetProjectCdsByProjectId(projectrecords[idx].ProjectId)

	if err != nil {

		return "", fmt.Errorf("failed to get cds: %s", err.Error())

	}

	cdlen := len(cds)

	if cdlen == 0 {

		return "", fmt.Errorf("failed to get cds: len: %d", cdlen)
	}

	clusterrecords, err := pkgdbquery.GetClustersByUserId(userrecord.UserId)

	if err != nil {

		return "", fmt.Errorf("failed to get clusters: %s", err.Error())
	}

	clen := len(clusterrecords)

	if clen == 0 {

		return "", fmt.Errorf("failed to get clusters: %d", clen)
	}

	out_struct := make([]pkgresourceapix.V1ProjectCdExport, 0)

	for i := 0; i < cdlen; i++ {

		tmpout := pkgresourceapix.V1ProjectCdExport{}

		tmpout.ProjectName = projectname

		cfound := -1

		for j := 0; j < clen; j++ {

			if cds[i].ClusterId == clusterrecords[j].ClusterId {

				cfound = j

				break

			}

		}

		if cfound == -1 {

			tmpout.ClusterName = pkgresourceapix.V1EXPORT_STRING_UNKNOWN

		} else {

			tmpout.ClusterName = clusterrecords[cfound].ClusterName
		}

		tmpout.ProjectCdStatus = cds[i].ProjectCdStatus

		if !cds[i].ProjectCdLog.Valid {

			tmpout.ProjectCdLog = pkgresourceapix.V1EXPORT_STRING_NULL

		} else {

			tmpout.ProjectCdLog = cds[i].ProjectCdLog.String

		}

		tmpout.ProjectCdStart = cds[i].ProjectCdStart

		if !cds[i].ProjectCdEnd.Valid {

			tmpout.ProjectCdEnd = time.Time{}

		} else {
			tmpout.ProjectCdEnd = cds[i].ProjectCdEnd.Time
		}

		out_struct = append(out_struct, tmpout)

	}

	outb, err := yaml.Marshal(out_struct)

	if err != nil {

		return "", fmt.Errorf("failed to marshal cd history: %s", err.Error())
	}

	output = string(outb)

	return output, nil
}

func V1LifecycleReportGetLatest(username string, projectname string) (string, error) {

	var output string

	userrecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return "", fmt.Errorf("failed to get lc latest: no such user: %s", err.Error())
	}

	if userrecord == nil {

		return "", fmt.Errorf("failed to get lc latest: empty user record")
	}

	projectrecords, err := pkgdbquery.GetProjectsByUserId(userrecord.UserId)

	if err != nil {

		return "", fmt.Errorf("failed to get lc latest: failed to get projects for: %s", username)
	}

	plen := len(projectrecords)

	if plen == 0 {

		return "", fmt.Errorf("failed to get lc latest: no associated projects for: %s", username)
	}

	idx := -1

	for i := 0; i < plen; i++ {

		if projectrecords[i].ProjectName == projectname {

			idx = i

			break

		}
	}

	if idx == -1 {

		return "", fmt.Errorf("failed to get lc latest: no such project: %s", projectname)
	}

	lcrecords, err := pkgdbquery.GetLifecyclesByProjectId(projectrecords[idx].ProjectId)

	if err != nil {

		return "", fmt.Errorf("failed to get lc latest: %s", err.Error())
	}

	lclen := len(lcrecords)

	if lclen == 0 {

		return "", fmt.Errorf("failed to get lc latest: len: %d", lclen)
	}

	out_struct := pkgresourceapix.V1LifecycleReportExport{}

	lc_latest := lcrecords[lclen-1]

	out_struct.ProjectName = projectname

	if !lc_latest.LifecycleReport.Valid {

		out_struct.LifecycleReport = pkgresourceapix.V1EXPORT_STRING_NULL

	} else {

		out_struct.LifecycleReport = lc_latest.LifecycleReport.String
	}

	out_struct.LifecycleStart = lc_latest.LifecycleStart

	outb, err := yaml.Marshal(out_struct)

	if err != nil {

		return "", fmt.Errorf("failed to get lc latest: marshal: %s", err.Error())
	}

	output = string(outb)

	return output, nil

}
