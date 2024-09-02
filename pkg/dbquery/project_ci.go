package dbquery

import (
	"fmt"

	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetProjectCi() ([]pkgresourcedb.DB_Project_CI, error) {

	var dbpci_records []pkgresourcedb.DB_Project_CI

	var dbpci pkgresourcedb.DB_Project_CI

	q := `
	
		SELECT
			project_ci_id,
			project_id,
			cluster_id,
			project_ci_status,
			project_ci_log,
			project_ci_start,
			project_ci_end
		FROM
			project_ci
	
	`

	a := []any{}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project ci: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpci = pkgresourcedb.DB_Project_CI{}

		err = res.Scan(
			&dbpci.ProjectCiId,
			&dbpci.ProjectId,
			&dbpci.ClusterId,
			&dbpci.ProjectCiStatus,
			&dbpci.ProjectCiLog,
			&dbpci.ProjectCiStart,
			&dbpci.ProjectCiEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project ci: row: %s", err.Error())
		}

		dbpci_records = append(dbpci_records, dbpci)
	}

	return dbpci_records, nil

}

func GetProjectCiById(id int) (*pkgresourcedb.DB_Project_CI, error) {

	var dbpci_records []pkgresourcedb.DB_Project_CI

	var dbpci pkgresourcedb.DB_Project_CI

	q := `
	
		SELECT
			project_id,
			cluster_id,
			project_ci_status,
			project_ci_log,
			project_ci_start,
			project_ci_end
		FROM
			project_ci
		WHERE
			project_ci_id = ?
	
	`

	a := []any{
		id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project ci by id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpci = pkgresourcedb.DB_Project_CI{}

		err = res.Scan(
			&dbpci.ProjectId,
			&dbpci.ClusterId,
			&dbpci.ProjectCiStatus,
			&dbpci.ProjectCiLog,
			&dbpci.ProjectCiStart,
			&dbpci.ProjectCiEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project ci by id: row: %s", err.Error())
		}

		dbpci_records = append(dbpci_records, dbpci)
	}

	rlen := len(dbpci_records)

	if rlen != 1 {
		return nil, fmt.Errorf("failed to get project ci by id: length: %d", rlen)
	}

	dbpci = dbpci_records[0]

	return &dbpci, nil

}

func GetProjectCisByProjectId(project_id int) ([]pkgresourcedb.DB_Project_CI, error) {

	var dbpci_records []pkgresourcedb.DB_Project_CI

	var dbpci pkgresourcedb.DB_Project_CI

	q := `
	
		SELECT
			project_ci_id,
			cluster_id,
			project_ci_status,
			project_ci_log,
			project_ci_start,
			project_ci_end
		FROM
			project_ci
		WHERE
			project_id = ?
	
	`

	a := []any{
		project_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cis by project id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpci = pkgresourcedb.DB_Project_CI{}

		err = res.Scan(
			&dbpci.ProjectCiId,
			&dbpci.ClusterId,
			&dbpci.ProjectCiStatus,
			&dbpci.ProjectCiLog,
			&dbpci.ProjectCiStart,
			&dbpci.ProjectCiEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cis by project id: row: %s", err.Error())
		}

		dbpci_records = append(dbpci_records, dbpci)
	}

	return dbpci_records, nil

}

func GetProjectCisByClusterId(cluster_id int) ([]pkgresourcedb.DB_Project_CI, error) {

	var dbpci_records []pkgresourcedb.DB_Project_CI

	var dbpci pkgresourcedb.DB_Project_CI

	q := `
	
		SELECT
			project_ci_id,
			project_id,
			project_ci_status,
			project_ci_log,
			project_ci_start,
			project_ci_end
		FROM
			project_ci
		WHERE
			cluster_id = ?
	
	`

	a := []any{
		cluster_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project cis by cluster id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpci = pkgresourcedb.DB_Project_CI{}

		err = res.Scan(
			&dbpci.ProjectCiId,
			&dbpci.ProjectId,
			&dbpci.ProjectCiStatus,
			&dbpci.ProjectCiLog,
			&dbpci.ProjectCiStart,
			&dbpci.ProjectCiEnd,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project cis by cluster id: row: %s", err.Error())
		}

		dbpci_records = append(dbpci_records, dbpci)
	}

	return dbpci_records, nil

}

func GetProjectCiRunningByProjectIdAndClusterId(project_id int, cluster_id int) (*pkgresourcedb.DB_Project_CI, error) {

	var dbpci_records []pkgresourcedb.DB_Project_CI

	var dbpci pkgresourcedb.DB_Project_CI

	q := `
	
		SELECT
			project_ci_id
		FROM
			project_ci
		WHERE
			project_id = ?
			AND cluster_id = ?
			AND project_ci_start IS NOT NULL
			AND project_ci_end IS NULL
	`

	a := []any{

		project_id,
		cluster_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project ci running by project id and cluster id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpci = pkgresourcedb.DB_Project_CI{}

		err = res.Scan(
			&dbpci.ProjectCiId,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project ci running by project id and cluster id: row: %s", err.Error())
		}

		dbpci_records = append(dbpci_records, dbpci)
	}

	rlen := len(dbpci_records)

	if rlen != 1 {

		return nil, fmt.Errorf("failed to get project ci running by project id and cluster id: len: %d", rlen)
	}

	dbpci = dbpci_records[0]

	return &dbpci, nil
}

func GetProjectCiRunningByProjectId(project_id int) (*pkgresourcedb.DB_Project_CI, error) {

	var dbpci_records []pkgresourcedb.DB_Project_CI

	var dbpci pkgresourcedb.DB_Project_CI

	q := `
	
		SELECT
			project_ci_id
		FROM
			project_ci
		WHERE
			project_id = ?
			AND project_ci_start IS NOT NULL
			AND project_ci_end IS NULL
	`

	a := []any{

		project_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get project ci running by project id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbpci = pkgresourcedb.DB_Project_CI{}

		err = res.Scan(
			&dbpci.ProjectCiId,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get project ci running by project id: row: %s", err.Error())
		}

		dbpci_records = append(dbpci_records, dbpci)
	}

	rlen := len(dbpci_records)

	if rlen != 1 {

		return nil, fmt.Errorf("failed to get project ci running by project id: len: %d", rlen)
	}

	dbpci = dbpci_records[0]

	return &dbpci, nil
}

func SetProjectCi(project_id int, cluster_id int) (*pkgresourcedb.DB_Project_CI, error) {

	q := `
	
		INSERT INTO 
			project_ci(
				project_id,
				cluster_id,
				project_ci_status,
				project_ci_log,
				project_ci_start,
				project_ci_end
			)
			VALUES(
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
		pkgresourceci.STATUS_READY,
	}

	err := DbExec(q, a)

	if err != nil {
		return nil, fmt.Errorf("failed to set project ci: %s", err.Error())
	}

	newpci, err := GetProjectCiRunningByProjectIdAndClusterId(project_id, cluster_id)

	if err != nil {

		return nil, fmt.Errorf("failed to set project ci: get new: %s", err.Error())
	}

	return newpci, nil
}

func SetProjectCiLogById(id int, ci_log string) error {

	check, err := GetProjectCiById(id)

	if err != nil {

		return fmt.Errorf("failed to set project ci log by id: %s", err.Error())
	}

	if check == nil {

		return fmt.Errorf("failed to set project ci log by id: not found: %d", id)
	}

	q := `

		UPDATE
			project_ci
		SET
			project_ci_status = ?,
			project_ci_log = ?
		WHERE
			project_ci_id = ?
	
	
	`

	a := []any{
		pkgresourceci.STATUS_RUNNING,
		ci_log,
		check.ProjectCiId,
	}

	err = DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set project ci log: %s", err.Error())
	}

	return nil
}

func SetProjectCiEndById(id int, ci_status string, ci_log string) error {

	check, err := GetProjectCiById(id)

	if err != nil {

		return fmt.Errorf("failed to set project ci end by id: %s", err.Error())
	}

	if check == nil {

		return fmt.Errorf("failed to set project ci end by id: not found: %d", id)
	}

	q := `

		UPDATE
			project_ci
		SET
			project_ci_status = ?,
			project_ci_log = ?,
			project_ci_end = CURRENT_TIMESTAMP(3)
		WHERE
			project_ci_id = ?
	
	
	`

	a := []any{
		ci_status,
		ci_log,
		check.ProjectCiId,
	}

	err = DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set project ci end: %s", err.Error())
	}

	return nil

}
