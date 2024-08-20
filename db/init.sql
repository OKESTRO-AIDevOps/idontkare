CREATE USER 'idkdb'@'%' IDENTIFIED BY 'universalpassword';

GRANT ALL PRIVILEGES ON *.* TO 'idkdb'@'%';

CREATE DATABASE idkdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE idkdb;


CREATE TABLE root (

    root_id INT NOT NULL AUTO_INCREMENT,
    root_name VARCHAR(128),
    root_ca_crt_path TEXT,
    root_ca_priv_path TEXT,
    root_server_crt_path TEXT,
    PRIMARY KEY (root_id)

);

CREATE TABLE user (

    user_id INT NOT NULL AUTO_INCREMENT,
    user_name VARCHAR(128),
    user_pass TEXT,
    PRIMARY KEY (user_id)

);

CREATE TABLE cluster (

    cluster_id INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    cluster_name VARCHAR(128),
    cluster_pub TEXT,
    cluster_connected TINYINT,
    cluster_session_key TEXT,
    PRIMARY KEY (cluster_id)

);


CREATE TABLE project (

    project_id INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    project_name VARCHAR(128),
    project_git TEXT,
    project_git_id TEXT,
    project_git_pw TEXT,
    project_registry TEXT,
    project_registry_id TEXT,
    project_registry_pw TEXT,
    project_cluster_id INT,
    project_ci_option TEXT,
    project_cd_option TEXT,
    PRIMARY KEY (project_id)

);


CREATE TABLE project_ci (

    project_ci_id INT NOT NULL AUTO_INCREMENT,
    project_id INT,
    cluster_id INT,
    project_ci_status VARCHAR(128),
    project_ci_log TEXT,
    project_ci_start TIMESTAMP(3),
    project_ci_end TIMESTAMP(3),
    PRIMARY KEY (project_ci_id)

);

CREATE TABLE project_cd (

    project_cd_id INT NOT NULL AUTO_INCREMENT,
    project_id INT,
    cluster_id INT,
    project_ci_id INT,
    project_cd_status VARCHAR(128),
    project_cd_log TEXT,
    project_cd_start TIMESTAMP(3),
    project_cd_end TIMESTAMP(3),
    PRIMARY KEY (project_cd_id)
);

COMMIT;