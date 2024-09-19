package dbquery

import (
	"fmt"

	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetLifecycle() ([]pkgresourcedb.DB_Lifecycle, error) {

	var dblifecycle_records []pkgresourcedb.DB_Lifecycle

	var dblifecycle pkgresourcedb.DB_Lifecycle

	q := `
	
		SELECT
			lifecycle_id,
			project_id,
			lifecycle_manifest,
			lifecycle_report,
			lifecycle_start
		FROM
			lifecycle
	
	
	`

	a := []any{}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get lifecycle: %s", err.Error())
	}

	defer res.Close()

	for res.Next() {

		dblifecycle = pkgresourcedb.DB_Lifecycle{}

		err = res.Scan(
			&dblifecycle.LifecycleId,
			&dblifecycle.ProjectId,
			&dblifecycle.LifecycleManifest,
			&dblifecycle.LifecycleReport,
			&dblifecycle.LifecycleStart,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get lifecycle: %s", err.Error())
		}

		dblifecycle_records = append(dblifecycle_records, dblifecycle)

	}

	return dblifecycle_records, nil

}

func GetLifecyclesByProjectId(project_id int) ([]pkgresourcedb.DB_Lifecycle, error) {

	var dblifecycle_records []pkgresourcedb.DB_Lifecycle

	var dblifecycle pkgresourcedb.DB_Lifecycle

	q := `
	
		SELECT
			lifecycle_id,
			lifecycle_manifest,
			lifecycle_report,
			lifecycle_start
		FROM
			lifecycle
		WHERE
			project_id = ?
	
	
	`

	a := []any{

		project_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get lifecycles by project id: %s", err.Error())
	}

	defer res.Close()

	for res.Next() {

		dblifecycle = pkgresourcedb.DB_Lifecycle{}

		err = res.Scan(
			&dblifecycle.LifecycleId,
			&dblifecycle.LifecycleManifest,
			&dblifecycle.LifecycleReport,
			&dblifecycle.LifecycleStart,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get lifecycles by project id: %s", err.Error())
		}

		dblifecycle_records = append(dblifecycle_records, dblifecycle)

	}

	return dblifecycle_records, nil
}

func GetLifecycleRunningByProjectId(project_id int) (*pkgresourcedb.DB_Lifecycle, error) {

	lcrecords, err := GetLifecyclesByProjectId(project_id)

	if err != nil {

		return nil, fmt.Errorf("failed to get lifecycle running by project id: %s", err.Error())
	}

	lclen := len(lcrecords)

	if lclen != 1 {

		return nil, fmt.Errorf("failed to get lifecycle running by project id: len: %s", err.Error())
	}

	lc := lcrecords[0]

	return &lc, nil
}

func SetLifecycleByProjectId(project_id int) error {

	q := `
	
	
		INSERT INTO
			lifecycle(
				project_id,
				lifecycle_manifest,
				lifecycle_report,
				lifecycle_start
			)
			VALUES(
				?,
				NULL,
				NULL,
				CURRENT_TIMESTAMP(3)
			)

	
	`

	a := []any{
		project_id,
	}

	err := DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set lifecycle by project id: %s", err.Error())
	}

	return nil

}

func SetLifecycleManifestByLifecycleId(lifecycle_id int, manifest string) error {

	q := `
	
		UPDATE
			lifecycle
		SET
			manifest = ?
		WHERE
			lifecycle_id = ?
	
	`

	a := []any{
		manifest,
		lifecycle_id,
	}

	err := DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set lifecycle manifest by id: %s", err.Error())
	}

	return nil
}

func SetLifecycleReportByLifecycleId(lifecycle_id int, report string) error {

	q := `
	
		UPDATE
			lifecycle
		SET
			report = ?
		WHERE
			lifecycle_id = ?
	
	`

	a := []any{
		report,
		lifecycle_id,
	}

	err := DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set lifecycle report by id: %s", err.Error())
	}

	return nil
}
