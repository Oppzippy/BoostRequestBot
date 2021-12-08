CREATE TABLE auto_signup_session (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	advertiser_id VARCHAR(50) NOT NULL,
	created_at DATETIME NOT NULL,
	expires_at DATETIME NOT NULL,
	deleted_at DATETIME,
	INDEX(guild_id, advertiser_id, expires_at, deleted_at),
	INDEX(expires_at, deleted_at)
);

