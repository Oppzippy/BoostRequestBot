CREATE TABLE log_channel (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	channel_id VARCHAR(50) NOT NULL,
	UNIQUE log_channel_guild_id (guild_id)
);
