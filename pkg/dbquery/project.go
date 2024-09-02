package dbquery

import (
	"fmt"
	"log"

	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetProject() ([]pkgresourcedb.DB_Project, error) {

	var dbproject_records []pkgresourcedb.DB_Project
	var dbproject pkgresourcedb.DB_Project

	q := `
		SELECT
		    project_id,
			user_id,
			project_name,
			project_git,
			project_git_id,
			project_git_pw,
			project_registry,
			project_registry_id,
			project_registry_pw,
			project_ci_option,
			project_cd_option
		FROM
			project
	
	`

	a := []any{}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbproject = pkgresourcedb.DB_Project{}

		err = res.Scan(
			&dbproject.ProjectId,
			&dbproject.UserId,
			&dbproject.ProjectName,
			&dbproject.ProjectGit,
			&dbproject.ProjectGitId,
			&dbproject.ProjectGitPw,
			&dbproject.ProjectRegistry,
			&dbproject.ProjectRegistryId,
			&dbproject.ProjectRegistryPw,
			&dbproject.ProjectCiOption,
			&dbproject.ProjectCdOption,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project by id: row: %s", err.Error())
		}

		dbproject_records = append(dbproject_records, dbproject)
	}

	return dbproject_records, nil
}

func GetProjectById(id int) (*pkgresourcedb.DB_Project, error) {

	var dbproject_records []pkgresourcedb.DB_Project

	var dbproject pkgresourcedb.DB_Project

	q := `
	
		SELECT
			user_id,
			project_name,
			project_git,
			project_git_id,
			project_git_pw,
			project_registry,
			project_registry_id,
			project_registry_pw,
			project_ci_option,
			project_cd_option
		FROM
			project
		WHERE
			project_id = ?
	
	`

	a := []any{
		id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project by id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbproject = pkgresourcedb.DB_Project{}

		err = res.Scan(
			&dbproject.UserId,
			&dbproject.ProjectName,
			&dbproject.ProjectGit,
			&dbproject.ProjectGitId,
			&dbproject.ProjectGitPw,
			&dbproject.ProjectRegistry,
			&dbproject.ProjectRegistryId,
			&dbproject.ProjectRegistryPw,
			&dbproject.ProjectCiOption,
			&dbproject.ProjectCdOption,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project by id: row: %s", err.Error())
		}

		dbproject_records = append(dbproject_records, dbproject)
	}

	rlen := len(dbproject_records)

	if rlen != 1 {
		return nil, nil
	}

	dbproject = dbproject_records[0]

	return &dbproject, nil

}

func GetProjectsByUserId(user_id int) ([]pkgresourcedb.DB_Project, error) {

	var dbproject_records []pkgresourcedb.DB_Project

	var dbproject pkgresourcedb.DB_Project

	q := `
	
		SELECT
			project_id,
			project_name,
			project_git,
			project_git_id,
			project_git_pw,
			project_registry,
			project_registry_id,
			project_registry_pw,
			project_ci_option,
			project_cd_option
		FROM
			project
		WHERE
			user_id = ?
	
	`

	a := []any{
		user_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project by user id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbproject = pkgresourcedb.DB_Project{}

		err = res.Scan(
			&dbproject.ProjectId,
			&dbproject.ProjectName,
			&dbproject.ProjectGit,
			&dbproject.ProjectGitId,
			&dbproject.ProjectGitPw,
			&dbproject.ProjectRegistry,
			&dbproject.ProjectRegistryId,
			&dbproject.ProjectRegistryPw,
			&dbproject.ProjectCiOption,
			&dbproject.ProjectCdOption,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project by user id: row: %s", err.Error())
		}

		dbproject_records = append(dbproject_records, dbproject)
	}

	return dbproject_records, nil

}

func SetProject(user_id int, name string, git string, gitid string, gitpw string, reg string, regid string, regpw string) error {

	exists := -1

	check, err := GetProjectsByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set project: %s", err.Error())

	}

	clen := len(check)

	for i := 0; i < clen; i++ {

		if check[i].ProjectName == name {

			exists = i

			break
		}

	}

	var q string

	a := make([]any, 0)

	if exists != -1 {

		log.Printf("project: %s: already exists: overriding\n", name)

		q = `
		
			UPDATE
				project
			SET
				user_id = ?,
				project_name = ?,
				project_git = ?,
				project_git_id = ?,
				project_git_pw = ?,
				project_registry = ?,
				project_registry_id = ?,
				project_registry_pw = ?,
				project_ci_option = NULL,
				project_cd_option = NULL
			WHERE
				project_id = ?

		
		`

		a = append(a, user_id)

		a = append(a, name)

		a = append(a, git)

		a = append(a, gitid)

		a = append(a, gitpw)

		a = append(a, reg)

		a = append(a, regid)

		a = append(a, regpw)

		a = append(a, check[exists].ProjectId)

	} else {

		log.Printf("project: %s: doesn't exist: create\n", name)

		q = `
		
			INSERT INTO
				project(
					user_id,
					project_name,
					project_git,
					project_git_id,
					project_git_pw,
					project_registry,
					project_registry_id,
					project_registry_pw,
					project_ci_option,
					project_cd_option
				)
				VALUES(
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
	                NULL,
	                NULL
				)

		
		`
		a = append(a, user_id)

		a = append(a, name)

		a = append(a, git)

		a = append(a, gitid)

		a = append(a, gitpw)

		a = append(a, reg)

		a = append(a, regid)

		a = append(a, regpw)
	}

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set cluster: %s", err.Error())
	}

	return nil

}

func SetProjectCiOptionByUserIdAndName(user_id int, name string, ci_option string) error {

	pexists := -1

	var q string

	a := make([]any, 0)

	pcheck, err := GetProjectsByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set project ci option: %s", err.Error())

	}

	pclen := len(pcheck)

	for i := 0; i < pclen; i++ {

		if pcheck[i].ProjectName == name {

			pexists = i

			break
		}

	}

	if pexists == -1 {

		return fmt.Errorf("failed to set project ci option: project doesn't exist: %s", name)

	}

	q = `
	
		UPDATE
			project
		SET
			project_ci_option = ?
		WHERE
			project_id = ?

	
	`

	a = append(a, ci_option)

	a = append(a, pcheck[pexists].ProjectId)

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set project ci option: %s", err.Error())
	}

	return nil

}
func SetProjectCdOptionByUserIdAndName(user_id int, name string, cd_option string) error {

	pexists := -1

	var q string

	a := make([]any, 0)

	pcheck, err := GetProjectsByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set project cd option: %s", err.Error())

	}

	pclen := len(pcheck)

	for i := 0; i < pclen; i++ {

		if pcheck[i].ProjectName == name {

			pexists = i

			break
		}

	}

	if pexists == -1 {

		return fmt.Errorf("failed to set project cd option: project doesn't exist: %s", name)

	}

	q = `
	
		UPDATE
			project
		SET
			project_cd_option = ?
		WHERE
			project_id = ?

	
	`
	a = append(a, cd_option)

	a = append(a, pcheck[pexists].ProjectId)

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set project cd option: %s", err.Error())
	}

	return nil

}
