package dbquery

import (
	"fmt"
	"log"

	pkgresourcedb "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/db"
)

func GetClusterById(id int) (*pkgresourcedb.DB_Cluster, error) {

	var dbcluster_records []pkgresourcedb.DB_Cluster

	var dbcluster pkgresourcedb.DB_Cluster

	q := `
	
		SELECT
			user_id,
			cluster_name,
			cluster_pub,
			cluster_connected,
			cluster_session_key
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
			cluster_session_key
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
