CREATE TABLE team_tb (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   name VARCHAR(15) NOT NULL UNIQUE,
   description TEXT NOT NULL
);

CREATE TABLE cluster_tb (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   description TEXT NOT NULL,
   name VARCHAR(25) NOT NULL UNIQUE,
   description TEXT NOT NULL,
   type VARCHAR(10) NOT NULL,
   address VARCHAR(64) NOT NULL UNIQUE,
   namespace VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE team_clusters(
   team_id INTEGER NOT NULL,
   cluster_id INTEGER NOT NULL
);

CREATE TABLE node_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   description TEXT NOT NULL,
   cluster_id INTEGER NOT NULL,
   name VARCHAR(25) NOT NULL UNIQUE,
   description TEXT NOT NULL,
   type VARCHAR(15) NOT NULL,
   address VARCHAR(64) NOT NULL,
   FOREIGN KEY(cluster_id) REFERENCES cluster_tb(id) ON DELETE CASCADE
);

CREATE TABLE application_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   name VARCHAR(25) NOT NULL UNIQUE,
   description TEXT NOT NULL,
   team_id INTEGER NOT NULL,
   cluster_id INTEGER NOT NULL,
   scale INTEGER NOT NULL,
   last_deployment TIMESTAMP,
   FOREIGN KEY(team_id) REFERENCES team_tb(id) ON DELETE CASCADE,
   FOREIGN KEY(cluster_id) REFERENCES cluster_tb(id) ON DELETE CASCADE
);

CREATE TABLE deployment_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   application_id INTEGER NOT NULL,
   image_uri VARCHAR(100) NOT NULL,
   report TEXT NULL,
   status VARCHAR(15) NOT NULL,
   added_at TIMESTAMP NOT NULL,
   completed_at TIMESTAMP NOT NULL,
   FOREIGN KEY(application_id) REFERENCES application_tb(id) ON DELETE CASCADE
);

CREATE TABLE container_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   ip VARCHAR(64) NOT NULL,
   family INTEGER NOT NULL,
   application_id INTEGER NOT NULL,
   node_id INTEGER NOT NULL,
   FOREIGN KEY(application_id) REFERENCES application_tb(id) ON DELETE CASCADE,
   FOREIGN KEY(node_id) REFERENCES node_tb(id) ON DELETE CASCADE,
);


CREATE TABLE environment_variable_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   key VARCHAR(50) NOT NULL,
   value TEXT NOT NULL,
   scope VARCHAR(50) NOT NULL,
   target VARCHAR(36) NOT NULL
);

CREATE TABLE user_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   nickname VARCHAR(36) NOT NULL,
   identifier VARCHAR(36) NOT NULL,
   email VARCHAR(100) NOT NULL,
   password VARCHAR(36) NOT NULL
);

CREATE TABLE user_teams(
   user_id INTEGER NOT NULL,
   team_id INTEGER NOT NULL
);

CREATE TABLE role_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name VARCHAR(50) NOT NULL
);

CREATE TABLE permission_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name VARCHAR(50) NOT NULL
);

CREATE TABLE role_permissions(
   role_id INTEGER NOT NULL,
   permission_id INTEGER NOT NULL
);

CREATE TABLE user_roles(
   user_id INTEGER NOT NULL,
   role_id INTEGER NOT NULL
);

CREATE TABLE session_tb(
   user_id INTEGER NOT NULL,
   identifier VARCHAR(36) NOT NULL,
   expiry TIMESTAMP NOT NULL
);

CREATE TRIGGER update_last_deployment AFTER UPDATE ON deployment_tb
	BEGIN
		UPDATE application_tb 
      SET last_deployment=(CASE WHEN NEW.status = 'DEPLOYED' THEN NEW.completed_at ELSE last_deployment END) 
      WHERE id=NEW.application_id;
	END;