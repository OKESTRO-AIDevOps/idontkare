package dbquery

import (
	"fmt"

	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetLifecyclesByProjectId(project_id int) ([]pkgresourcedb.DB_Lifecycle, error) {

	var dblifecycle_records []pkgresourcedb.DB_Lifecycle

	var dblifecycle pkgresourcedb.DB_Lifecycle

	q := `
	
		SELECT
			lifecycle_id,
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

func SetLifecycleByProjectId(project_id int) error {

	q := `
	
	
		INSERT INTO
			lifecycle(
				project_id,
				lifecycle_report,
				lifecycle_start
			)
			VALUES(
				?,
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
