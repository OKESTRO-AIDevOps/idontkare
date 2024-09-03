package dbquery

import (
	"database/sql"
	"fmt"
	"time"

	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DbEstablish(db_id string, db_pw string, db_addr string, db_name string) error {

	var err error

	db_info := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", db_id, db_pw, db_addr, db_name)

	DB, err = sql.Open("mysql", db_info)

	if err != nil {

		return fmt.Errorf("failed to open db: %s", err.Error())
	}

	DB.SetConnMaxLifetime(time.Second * 10)
	DB.SetConnMaxIdleTime(time.Second * 5)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	return nil
}

func DbQuery(query string, args []any) (*sql.Rows, error) {

	var empty_row *sql.Rows

	results, err := DB.Query(query, args[0:]...)

	if err != nil {

		return empty_row, fmt.Errorf("db query: %s", err.Error())

	}

	return results, err

}

func DbExec(query string, args []any) error {

	result, err := DB.Exec(query, args[0:]...)

	if err != nil {
		return fmt.Errorf("db exec: %s", err.Error())
	}

	_, err = result.RowsAffected()

	if err != nil {

		return fmt.Errorf("db exec: rows: %s", err.Error())
	}

	return nil

}

func DBReset() error {

	crecords, err := GetCluster()

	if err != nil {

		return fmt.Errorf("failed to reset: get clusters: %s", err.Error())
	}

	q := ""
	a := []any{}

	clen := len(crecords)

	for i := 0; i < clen; i++ {

		if crecords[i].ClusterConnected == 1 {

			q = `
	
			UPDATE
				cluster 
			SET
				cluster_connected = 0
			WHERE
				cluster_connected = 1
		
			`

			a = []any{}

			err = DbExec(q, a)

			if err != nil {

				return fmt.Errorf("failed to reset: cluster: %s", err.Error())
			}

			break
		}
	}

	pcdrecords, err := GetProjectCd()

	if err != nil {

		return fmt.Errorf("failed to reset: project cd: %s", err.Error())
	}

	pcdlen := len(pcdrecords)

	for i := 0; i < pcdlen; i++ {

		if !pcdrecords[i].ProjectCdEnd.Valid {

			q = `
			
			UPDATE
				project_cd
			SET
				project_cd_status = ?,
				project_cd_log = ?,
				project_cd_end = CURRENT_TIMESTAMP(3)

			WHERE
				project_cd_end IS NULL


			`

			a = []any{
				pkgresourcecd.STATUS_ERROR,
				"reset due to server policy",
			}

			err = DbExec(q, a)

			if err != nil {

				return fmt.Errorf("failed to reset: project cd: %s", err.Error())
			}

			break

		}

	}

	pcirecords, err := GetProjectCi()

	if err != nil {

		return fmt.Errorf("failed to reset: project ci: %s", err.Error())
	}

	pcilen := len(pcirecords)

	for i := 0; i < pcilen; i++ {

		if !pcirecords[i].ProjectCiEnd.Valid {

			q = `
	
			UPDATE
				project_ci
			SET
				project_ci_status = ?,
				project_ci_log = ?,
				project_ci_end = CURRENT_TIMESTAMP(3)
		
			WHERE
				project_ci_end IS NULL
		
		
			`

			a = []any{
				pkgresourceci.STATUS_ERROR,
				"reset due to server policy",
			}

			err = DbExec(q, a)

			if err != nil {

				return fmt.Errorf("failed to reset: project ci: %s", err.Error())
			}

			break

		}

	}

	return nil
}
