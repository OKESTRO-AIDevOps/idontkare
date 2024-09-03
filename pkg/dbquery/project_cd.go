package dbquery

import (
	"fmt"

	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetProjectCd() ([]pkgresourcedb.DB_Project_CD, error) {

	var dbpcd_records []pkgresourcedb.DB_Project_CD

	var dbpcd pkgresourcedb.DB_Project_CD

	q := `
	
		SELECT
			project_cd_id,
			project_id,
			cluster_id,
			project_ci_id,
			project_cd_status,
			project_cd_log,
			project_cd_start,
			project_cd_end
		FROM
			project_cd

	
	`

	a := []any{}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cd: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpcd = pkgresourcedb.DB_Project_CD{}

		err = res.Scan(
			&dbpcd.ProjectCdId,
			&dbpcd.ProjectId,
			&dbpcd.ClusterId,
			&dbpcd.ProjectCiId,
			&dbpcd.ProjectCdStatus,
			&dbpcd.ProjectCdLog,
			&dbpcd.ProjectCdStart,
			&dbpcd.ProjectCdEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cd: row: %s", err.Error())
		}

		dbpcd_records = append(dbpcd_records, dbpcd)
	}

	return dbpcd_records, nil
}

func GetProjectCdById(id int) (*pkgresourcedb.DB_Project_CD, error) {

	var dbpcd_records []pkgresourcedb.DB_Project_CD

	var dbpcd pkgresourcedb.DB_Project_CD

	q := `
	
		SELECT
			project_id,
			cluster_id,
			project_ci_id,
			project_cd_status,
			project_cd_log,
			project_cd_start,
			project_cd_end
		FROM
			project_cd
		WHERE
			project_cd_id = ?
	
	`

	a := []any{
		id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cd by id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpcd = pkgresourcedb.DB_Project_CD{}

		err = res.Scan(
			&dbpcd.ProjectId,
			&dbpcd.ClusterId,
			&dbpcd.ProjectCiId,
			&dbpcd.ProjectCdStatus,
			&dbpcd.ProjectCdLog,
			&dbpcd.ProjectCdStart,
			&dbpcd.ProjectCdEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cd by id: row: %s", err.Error())
		}

		dbpcd_records = append(dbpcd_records, dbpcd)
	}

	rlen := len(dbpcd_records)

	if rlen != 1 {
		return nil, fmt.Errorf("failed to get project cd by id: length: %d", rlen)
	}

	dbpcd = dbpcd_records[0]

	return &dbpcd, nil
}

func GetProjectCdsByProjectId(project_id int) ([]pkgresourcedb.DB_Project_CD, error) {

	var dbpcd_records []pkgresourcedb.DB_Project_CD

	var dbpcd pkgresourcedb.DB_Project_CD

	q := `
	
		SELECT
			project_cd_id,
			cluster_id,
			project_ci_id,
			project_cd_status,
			project_cd_log,
			project_cd_start,
			project_cd_end
		FROM
			project_cd
		WHERE
			project_id = ?
	
	`

	a := []any{
		project_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cds by project id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpcd = pkgresourcedb.DB_Project_CD{}

		err = res.Scan(
			&dbpcd.ProjectCdId,
			&dbpcd.ClusterId,
			&dbpcd.ProjectCiId,
			&dbpcd.ProjectCdStatus,
			&dbpcd.ProjectCdLog,
			&dbpcd.ProjectCdStart,
			&dbpcd.ProjectCdEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cds by project id: row: %s", err.Error())
		}

		dbpcd_records = append(dbpcd_records, dbpcd)
	}

	return dbpcd_records, nil

}

func GetProjectCdsByClusterId(cluster_id int) ([]pkgresourcedb.DB_Project_CD, error) {

	var dbpcd_records []pkgresourcedb.DB_Project_CD

	var dbpcd pkgresourcedb.DB_Project_CD

	q := `
	
		SELECT
			project_cd_id,
			project_id,
			project_ci_id,
			project_cd_status,
			project_cd_log,
			project_cd_start,
			project_cd_end
		FROM
			project_cd
		WHERE
			cluster_id = ?
	
	`

	a := []any{
		cluster_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cds by cluster id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpcd = pkgresourcedb.DB_Project_CD{}

		err = res.Scan(
			&dbpcd.ProjectCdId,
			&dbpcd.ProjectId,
			&dbpcd.ProjectCiId,
			&dbpcd.ProjectCdStatus,
			&dbpcd.ProjectCdLog,
			&dbpcd.ProjectCdStart,
			&dbpcd.ProjectCdEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cds by cluster id: row: %s", err.Error())
		}

		dbpcd_records = append(dbpcd_records, dbpcd)
	}

	return dbpcd_records, nil

}

func GetProjectCdsByCiId(ci_id int) ([]pkgresourcedb.DB_Project_CD, error) {

	var dbpcd_records []pkgresourcedb.DB_Project_CD

	var dbpcd pkgresourcedb.DB_Project_CD

	q := `
	
		SELECT
			project_cd_id,
			project_id,
			cluster_id,
			project_cd_status,
			project_cd_log,
			project_cd_start,
			project_cd_end
		FROM
			project_ci
		WHERE
			project_ci_id = ?
	
	`

	a := []any{
		ci_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cds by project ci id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpcd = pkgresourcedb.DB_Project_CD{}

		err = res.Scan(
			&dbpcd.ProjectCdId,
			&dbpcd.ProjectId,
			&dbpcd.ClusterId,
			&dbpcd.ProjectCdStatus,
			&dbpcd.ProjectCdLog,
			&dbpcd.ProjectCdStart,
			&dbpcd.ProjectCdEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cds by project ci id: row: %s", err.Error())
		}

		dbpcd_records = append(dbpcd_records, dbpcd)
	}

	return dbpcd_records, nil
}

func GetProjectCdRunningByProjectIdAndClusterId(project_id int, cluster_id int) (*pkgresourcedb.DB_Project_CD, error) {

	var dbpcd_records []pkgresourcedb.DB_Project_CD

	var dbpcd pkgresourcedb.DB_Project_CD

	q := `
	
		SELECT
			project_cd_id
		FROM
			project_cd
		WHERE
			project_id = ?
			AND cluster_id = ?
			AND project_cd_start IS NOT NULL
			AND project_cd_end IS NULL
	`

	a := []any{

		project_id,
		cluster_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cd running by project id and cluster id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpcd = pkgresourcedb.DB_Project_CD{}

		err = res.Scan(
			&dbpcd.ProjectCdId,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cd running by project id and cluster id: row: %s", err.Error())
		}

		dbpcd_records = append(dbpcd_records, dbpcd)
	}

	rlen := len(dbpcd_records)

	if rlen != 1 {

		return nil, fmt.Errorf("failed to get project cd running by project id and cluster id: len: %d", rlen)
	}

	dbpcd = dbpcd_records[0]

	return &dbpcd, nil
}

func SetProjectCd(project_id int, cluster_id int, ci_id int) (*pkgresourcedb.DB_Project_CD, error) {

	q := `
	
		INSERT INTO 
			project_cd(
				project_id,
				cluster_id,
				project_ci_id,
				project_cd_status,
				project_cd_log,
				project_cd_start,
				project_cd_end
			)
			VALUES(
				?,
				?,
				?,
				?,
				NULL,
				CURRENT_TIMESTAMP(3),
				NULL
			)
	
	`

	a := []any{
		project_id,
		cluster_id,
		ci_id,
		pkgresourcecd.STATUS_READY,
	}

	err := DbExec(q, a)

	if err != nil {
		return nil, fmt.Errorf("failed to set project cd: %s", err.Error())
	}

	newpcd, err := GetProjectCdRunningByProjectIdAndClusterId(project_id, cluster_id)

	if err != nil {

		return nil, fmt.Errorf("failed to set project cd: get new: %s", err.Error())
	}

	return newpcd, nil

}

func SetProjectCdLogById(id int, cd_log string) error {

	check, err := GetProjectCdById(id)

	if err != nil {

		return fmt.Errorf("failed to set project cd log by id: %s", err.Error())
	}

	if check == nil {

		return fmt.Errorf("failed to set project cd log by id: not found: %d", id)
	}

	q := `

		UPDATE
			project_cd
		SET
			project_cd_status = ?,
			project_cd_log = ?
		WHERE
			project_cd_id = ?
	
	
	`

	a := []any{
		pkgresourcecd.STATUS_RUNNING,
		cd_log,
		id,
	}

	err = DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set project cd log: %s", err.Error())
	}

	return nil
}

func SetProjectCdEndById(id int, cd_status string, cd_log string) error {

	check, err := GetProjectCdById(id)

	if err != nil {

		return fmt.Errorf("failed to set project cd end by id: %s", err.Error())
	}

	if check == nil {

		return fmt.Errorf("failed to set project cd end by id: not found: %d", id)
	}

	q := `

		UPDATE
			project_cd
		SET
			project_cd_status = ?,
			project_cd_log = ?,
			project_cd_end = CURRENT_TIMESTAMP(3)
		WHERE
			project_cd_id = ?
	
	
	`

	a := []any{
		cd_status,
		cd_log,
		id,
	}

	err = DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set project cd end: %s", err.Error())
	}

	return nil

}
