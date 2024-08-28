package apiximpl

import (
	"fmt"
	"strings"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
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
