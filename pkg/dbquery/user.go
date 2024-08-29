package dbquery

import (
	"fmt"
	"log"

	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetUserById(id int) (*pkgresourcedb.DB_User, error) {

	var dbuser_records []pkgresourcedb.DB_User

	var dbuser pkgresourcedb.DB_User

	q := `
	
		SELECT
			user_name,
			user_pass
		FROM
			user
		WHERE
			user_id = ?
	
	`

	a := []any{
		id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get user: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbuser := pkgresourcedb.DB_User{}

		err = res.Scan(
			&dbuser.UserName,
			&dbuser.UserPass,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get user: row: %s", err.Error())
		}

		dbuser_records = append(dbuser_records, dbuser)
	}

	rlen := len(dbuser_records)

	if rlen != 1 {
		return nil, nil
	}

	dbuser = dbuser_records[0]

	return &dbuser, nil

}

func GetUserByName(name string) (*pkgresourcedb.DB_User, error) {

	var dbuser_records []pkgresourcedb.DB_User

	var dbuser pkgresourcedb.DB_User

	q := `
	
		SELECT
			user_id,
			user_pass
		FROM
			user
		WHERE
			user_name = ?
	
	`

	a := []any{
		name,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get user: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbuser := pkgresourcedb.DB_User{}

		err = res.Scan(
			&dbuser.UserId,
			&dbuser.UserPass,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get user: row: %s", err.Error())
		}

		dbuser_records = append(dbuser_records, dbuser)
	}

	rlen := len(dbuser_records)

	if rlen != 1 {
		return nil, nil
	}

	dbuser = dbuser_records[0]

	return &dbuser, nil
}

func SetUser(name string, pass string) error {

	check, err := GetUserByName(name)

	if err != nil {

		return fmt.Errorf("set user: %s", err.Error())
	}

	var q string

	a := make([]any, 0)

	if check != nil {

		log.Printf("user: %s: already exists: overriding\n", name)

		q = `
		
			UPDATE
				user
			SET
				user_pass = ?
			WHERE
				user_id = ?

		
		`

		a = append(a, pass)

		a = append(a, check.UserId)

	} else {

		log.Printf("user: %s: doesn't exist: create\n", name)

		q = `
		
			INSERT INTO
				user(
					user_name,
					user_pass
				)
				VALUES(
					?,
					?
				)

		
		`

		a = append(a, name)
		a = append(a, pass)
	}

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set user: %s", err.Error())
	}

	return nil
}
