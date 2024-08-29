package dbquery

import (
	"fmt"

	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetRoot() (*pkgresourcedb.DB_Root, error) {

	var dbroot_records []pkgresourcedb.DB_Root

	var dbroot pkgresourcedb.DB_Root

	q := `
	
		SELECT
			root_id,
			root_name,
			root_ca_crt_path,
			root_ca_priv_path,
			root_server_crt_path
		FROM
			root

	
	`

	a := []any{}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get root info: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbroot = pkgresourcedb.DB_Root{}

		err = res.Scan(
			&dbroot.RootId,
			&dbroot.RootName,
			&dbroot.RootCACrtPath,
			&dbroot.RootCAPrivPath,
			&dbroot.RootServerCrtPath,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get root info: row: %s", err.Error())

		}

		dbroot_records = append(dbroot_records, dbroot)
	}

	rlen := len(dbroot_records)

	if rlen != 1 {
		return nil, fmt.Errorf("failed to get root info: length: %d", rlen)
	}

	dbroot = dbroot_records[0]

	return &dbroot, nil
}
