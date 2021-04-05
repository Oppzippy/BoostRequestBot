CREATE TABLE advertiser_privileges (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	role_id VARCHAR(50) NOT NULL,
	weight DOUBLE PRECISION NOT NULL,
	delay INTEGER NOT NULL,
	created_at DATETIME NOT NULL,
	UNIQUE advertiser_privileges_unique (guild_id, role_id)
);
