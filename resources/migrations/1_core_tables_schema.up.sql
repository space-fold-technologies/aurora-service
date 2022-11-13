CREATE TABLE team_tb (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   name VARCHAR(50) NOT NULL UNIQUE,
   description TEXT NOT NULL
);

CREATE TABLE cluster_tb (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   description TEXT NOT NULL,
   name VARCHAR(25) NOT NULL UNIQUE,
   type VARCHAR(10) NOT NULL,
   address VARCHAR(64) NOT NULL UNIQUE,
   token VARCHAR(256) NOT NULL UNIQUE,
   namespace VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE cluster_teams(
   cluster_id INTEGER NOT NULL,
   team_id INTEGER NOT NULL
);

CREATE TABLE node_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   description TEXT NOT NULL,
   cluster_id INTEGER NOT NULL,
   name VARCHAR(25) NOT NULL UNIQUE,
   type VARCHAR(15) NOT NULL,
   address VARCHAR(64) NOT NULL,
   FOREIGN KEY(cluster_id) REFERENCES cluster_tb(id) ON DELETE CASCADE
);

CREATE TABLE deployment_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   application_id INTEGER NOT NULL,
   image_uri VARCHAR(100) NOT NULL,
   report TEXT NULL,
   service_identifier VARCHAR(36) NULL,
   status VARCHAR(255) NOT NULL,
   added_at TIMESTAMP NOT NULL,
   completed_at TIMESTAMP,
   FOREIGN KEY(application_id) REFERENCES application_tb(id) ON DELETE CASCADE
);


CREATE TABLE container_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   identifier VARCHAR(36) NOT NULL,
   ip VARCHAR(64) NOT NULL,
   address_family INTEGER NOT NULL,
   application_id INTEGER NOT NULL,
   node_id INTEGER NOT NULL,
   FOREIGN KEY(application_id) REFERENCES application_tb(id) ON DELETE CASCADE,
   FOREIGN KEY(node_id) REFERENCES node_tb(id) ON DELETE CASCADE
);

CREATE TABLE environment_variable_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   key VARCHAR(50) NOT NULL,
   value TEXT NOT NULL,
   scope VARCHAR(50) NOT NULL,
   target VARCHAR(36) NOT NULL
);

CREATE TABLE permission_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name VARCHAR(50) NOT NULL
);

CREATE TABLE user_tb(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name VARCHAR(100) NOT NULL,
   nick_name VARCHAR(36) NOT NULL,
   identifier VARCHAR(36) NOT NULL,
   email VARCHAR(100) NOT NULL,
   password VARCHAR(36) NOT NULL
);

CREATE TABLE user_teams(
   user_id INTEGER NOT NULL,
   team_id INTEGER NOT NULL
);


CREATE TABLE user_permissions(
   user_id INTEGER NOT NULL,
   permission_id INTEGER NOT NULL
);

CREATE TRIGGER update_last_deployment AFTER UPDATE ON deployment_tb
	BEGIN
		UPDATE application_tb 
      SET last_deployment=(CASE WHEN NEW.status = 'DEPLOYED' THEN NEW.completed_at ELSE last_deployment END) 
      WHERE id=NEW.application_id;
	END;