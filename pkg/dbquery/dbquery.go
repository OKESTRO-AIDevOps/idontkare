package dbquery

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DbEstablish(db_id string, db_pw string, db_addr string, db_name string) error {

	var err error

	db_info := fmt.Sprintf("%s:%s@tcp(%s)/%s", db_id, db_pw, db_addr, db_name)

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

	count, err := result.RowsAffected()

	if err != nil {

		return fmt.Errorf("db exec: rows: %s", err.Error())
	}

	if count < 1 {
		return fmt.Errorf("db exec: 0 affected")
	}

	return nil

}

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

func GetClusterById(id int) (*pkgresourcedb.DB_Cluster, error) {

	var dbcluster_records []pkgresourcedb.DB_Cluster

	var dbcluster pkgresourcedb.DB_Cluster

	q := `
	
		SELECT
			user_id,
			cluster_name,
			cluster_pub,
			cluster_connected,
			cluster_session_key,
		FROM
			cluster
		WHERE
			cluster_id = ?
	
	`

	a := []any{
		id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get cluster by id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbcluster = pkgresourcedb.DB_Cluster{}

		err = res.Scan(
			&dbcluster.UserId,
			&dbcluster.ClusterName,
			&dbcluster.ClusterPub,
			&dbcluster.ClusterConnected,
			&dbcluster.ClusterSessionKey,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get cluster by id: row: %s", err.Error())
		}

		dbcluster_records = append(dbcluster_records, dbcluster)
	}

	rlen := len(dbcluster_records)

	if rlen != 1 {
		return nil, nil
	}

	dbcluster = dbcluster_records[0]

	return &dbcluster, nil
}

func GetClustersByUserId(user_id int) ([]pkgresourcedb.DB_Cluster, error) {

	var dbcluster_records []pkgresourcedb.DB_Cluster

	var dbcluster pkgresourcedb.DB_Cluster

	q := `
	
		SELECT
			cluster_id,
			cluster_name,
			cluster_pub,
			cluster_connected,
			cluster_session_key,
		FROM
			cluster
		WHERE
			user_id = ?
	
	`

	a := []any{
		user_id,
	}

	res, err := DbQuery(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get cluster by user id: %s", err.Error())

	}

	defer res.Close()

	for res.Next() {

		dbcluster = pkgresourcedb.DB_Cluster{}

		err = res.Scan(
			&dbcluster.ClusterId,
			&dbcluster.ClusterName,
			&dbcluster.ClusterPub,
			&dbcluster.ClusterConnected,
			&dbcluster.ClusterSessionKey,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get cluster by user id: row: %s", err.Error())
		}

		dbcluster_records = append(dbcluster_records, dbcluster)
	}

	return dbcluster_records, nil
}

func SetCluster(user_id int, name string, pub string) error {

	exists := -1

	check, err := GetClustersByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set cluster: %s", err.Error())

	}

	clen := len(check)

	for i := 0; i < clen; i++ {

		if check[i].ClusterName == name {

			exists = i

			break
		}

	}

	var q string

	a := make([]any, 0)

	if exists != -1 {

		log.Printf("cluster: %s: already exists: overriding\n", name)

		q = `
		
			UPDATE
				cluster
			SET
				user_id = ?,
				cluster_name = ?,
				cluster_pub = ?,
				cluster_connected = 0,
				cluster_session_key = NULL
			WHERE
				cluster_id = ?

		
		`

		a = append(a, user_id)

		a = append(a, name)

		a = append(a, pub)

		a = append(a, check[exists].ClusterId)

	} else {

		log.Printf("cluster: %s: doesn't exist: create\n", name)

		q = `
		
			INSERT INTO
				cluster(
					user_id,
					cluster_name,
					cluster_pub,
					cluster_connected,
					cluster_session_key
				)
				VALUES(
					?,
					?,
					?,
					0,
					NULL
				)

		
		`
		a = append(a, user_id)

		a = append(a, name)

		a = append(a, pub)
	}

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set cluster: %s", err.Error())
	}

	return nil
}

func SetClusterConnectedByUserIdAndName(user_id int, name string, connected int, session_key string) error {

	exists := -1

	check, err := GetClustersByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set cluster connected: %s", err.Error())

	}

	clen := len(check)

	for i := 0; i < clen; i++ {

		if check[i].ClusterName == name {

			exists = i

			break
		}

	}

	var q string

	a := make([]any, 0)

	if exists == -1 {

		return fmt.Errorf("failed to set cluster connected: doesn't exist: %s", name)

	} else {

		log.Printf("cluster: %s: change connected to: %d: session key: %s\n", name, connected, session_key)

		q = `
		
			UPDATE
				cluster
			SET
				cluster_connected = ?,
				cluster_session_key = ?
			WHERE
				cluster_id = ?

		
		`
		a = append(a, connected)

		a = append(a, session_key)

		a = append(a, check[exists].ClusterId)
	}

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set cluster connected: %s", err.Error())
	}

	return nil

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
			project_cluster_id,
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
			&dbproject.ProjectClusterId,
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
			project_cluster_id,
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
			&dbproject.UserId,
			&dbproject.ProjectName,
			&dbproject.ProjectGit,
			&dbproject.ProjectGitId,
			&dbproject.ProjectGitPw,
			&dbproject.ProjectRegistry,
			&dbproject.ProjectRegistryId,
			&dbproject.ProjectRegistryPw,
			&dbproject.ProjectClusterId,
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
				project_cluster_id = -1,
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
					project_cluster_id,
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
					-1,
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

func SetProjectClusterIdByUserIdAndName(user_id int, name string, cluster_id int) error {

	pexists := -1
	cexists := -1

	var q string

	a := make([]any, 0)

	pcheck, err := GetProjectsByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set project cluster id: %s", err.Error())

	}

	ccheck, err := GetClustersByUserId(user_id)

	if err != nil {

		return fmt.Errorf("set project cluster id: %s", err.Error())

	}

	pclen := len(pcheck)

	cclen := len(ccheck)

	for i := 0; i < pclen; i++ {

		if pcheck[i].ProjectName == name {

			pexists = i

			break
		}

	}

	for i := 0; i < cclen; i++ {

		if ccheck[i].ClusterId == cluster_id {

			cexists = i

			break
		}
	}

	if pexists == -1 {

		return fmt.Errorf("failed to set project cluster id: project doesn't exist: %s", name)

	}

	if cexists == -1 {

		return fmt.Errorf("failed to set project cluster id: cluster doesn't exist: %s", name)

	}

	q = `
	
		UPDATE
			project
		SET
			project_cluster_id = ?
		WHERE
			project_id = ?

	
	`
	a = append(a, cluster_id)

	a = append(a, pcheck[pexists].ProjectId)

	err = DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set project cluster id: %s", err.Error())
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
			cluster_id,
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

func SetProjectCi(project_id int, cluster_id int) error {

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
		return fmt.Errorf("failed to set project ci: %s", err.Error())
	}

	return nil
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

func SetProjectCd(project_id int, cluster_id int, ci_id int) error {

	q := `
	
		INSERT INTO 
			project_cd(
				project_id,
				cluster_id,
				project_ci_id,
				project_ci_status,
				project_ci_log,
				project_ci_start,
				project_ci_end
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
		pkgresourceci.STATUS_READY,
	}

	err := DbExec(q, a)

	if err != nil {
		return fmt.Errorf("failed to set project cd: %s", err.Error())
	}

	return nil

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
		pkgresourceci.STATUS_RUNNING,
		cd_log,
		check.ProjectCdId,
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
		check.ProjectCdId,
	}

	err = DbExec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set project cd end: %s", err.Error())
	}

	return nil

}
