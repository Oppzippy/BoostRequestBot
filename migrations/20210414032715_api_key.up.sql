CREATE TABLE api_key (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	`key` VARCHAR(50) NOT NULL,
	guild_id VARCHAR(50) NOT NULL,
	created_at DATETIME NOT NULL,
	UNIQUE api_key_unique (`key`)
);
