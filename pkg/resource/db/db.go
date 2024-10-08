package resource

import (
	"database/sql"
	"time"
)

type DB_Root struct {
	RootId            int    `json:"root_id"`
	RootName          string `json:"root_name"`
	RootCACrtPath     string `json:"root_ca_crt_path"`
	RootCAPrivPath    string `json:"root_ca_priv_path"`
	RootServerCrtPath string `json:"root_server_crt_path"`
}

type DB_User struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	UserPass string `json:"user_pass"`
}

type DB_Cluster struct {
	ClusterId         int            `json:"cluster_id"`
	UserId            int            `json:"user_id"`
	ClusterName       string         `json:"cluster_name"`
	ClusterPub        string         `json:"cluster_pub"`
	ClusterConnected  int            `json:"cluster_connected"`
	ClusterSessionKey sql.NullString `json:"cluster_session_key"`
}

type DB_Project struct {
	ProjectId         int            `json:"project_id"`
	UserId            int            `json:"user_id"`
	ProjectName       string         `json:"project_name"`
	ProjectGit        string         `json:"project_git"`
	ProjectGitId      string         `json:"project_git_id"`
	ProjectGitPw      string         `json:"project_git_pw"`
	ProjectRegistry   string         `json:"project_registry"`
	ProjectRegistryId string         `json:"project_registry_id"`
	ProjectRegistryPw string         `json:"project_registry_pw"`
	ProjectCiOption   sql.NullString `json:"project_ci_option"`
	ProjectCdOption   sql.NullString `json:"project_cd_option"`
}

type DB_Project_CI struct {
	ProjectCiId     int            `json:"project_ci_id"`
	ProjectId       int            `json:"project_id"`
	ClusterId       int            `json:"cluster_id"`
	ProjectCiStatus string         `json:"project_ci_status"`
	ProjectCiLog    sql.NullString `json:"project_ci_log"`
	ProjectCiStart  time.Time      `json:"project_ci_start"`
	ProjectCiEnd    sql.NullTime   `json:"project_ci_end"`
}

type DB_Project_CD struct {
	ProjectCdId     int            `json:"project_cd_id"`
	ProjectId       int            `json:"project_id"`
	ProjectCiId     int            `json:"project_ci_id"`
	ClusterId       int            `json:"cluster_id"`
	ProjectCdStatus string         `json:"project_cd_status"`
	ProjectCdLog    sql.NullString `json:"project_cd_log"`
	ProjectCdStart  time.Time      `json:"project_cd_start"`
	ProjectCdEnd    sql.NullTime   `json:"project_cd_end"`
}

type DB_Lifecycle struct {
	LifecycleId       int            `json:"lifecycle_id"`
	ProjectId         int            `json:"project_id"`
	LifecycleManifest sql.NullString `json:"lifecycle_manifest"`
	LifecycleReport   sql.NullString `json:"lifecycle_report"`
	LifecycleStart    time.Time      `json:"lifecycle_start"`
}
