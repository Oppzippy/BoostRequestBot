CREATE TABLE boost_request_channel (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	frontend_channel_id VARCHAR(50) NOT NULL,
	backend_channel_id VARCHAR(50) NOT NULL,
	uses_buyer_message TINYINT(1) NOT NULL,
	skips_buyer_dm TINYINT(1) NOT NULL,
	created_at DATETIME NOT NULL,
	UNIQUE boost_request_channel_unique (guild_id, frontend_channel_id),
	INDEX boost_request_channel_backend_channel_index (backend_channel_id)
);

CREATE TABLE boost_request (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_channel_id INTEGER NOT NULL,
	requester_id VARCHAR(50) NOT NULL,
	advertiser_id VARCHAR(50) NULL,
	backend_message_id VARCHAR(50) NOT NULL,
	message TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	resolved_at DATETIME NULL,
	deleted_at DATETIME NULL,
	UNIQUE boost_request_unique (backend_message_id)
);
