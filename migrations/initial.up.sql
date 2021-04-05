CREATE TABLE boost_request_channel (
	id INTEGER PRIMARY KEY AUTO_INCREMENT
	guild_id VARCHAR(50) NOT NULL,
	frontend_channel_id VARCHAR(50) NOT NULL,
	backend_channel_id VARCHAR(50) NOT NULL,
	uses_buyer_message TINYINT(1) NOT NULL,
	notifies_buyer TINYINT(1) NOT NULL,
	created_at DATETIME NOT NULL,
	deleted_at DATETIME NULL
);

CREATE INDEX boost_request_channel_frontend_channel_id_index USING BTREE ON boost_request_bot.boost_request_channel (frontend_channel_id);
CREATE INDEX boost_request_channel_backend_channel_id_index USING BTREE ON boost_request_bot.boost_request_channel (backend_channel_id);


CREATE TABLE boost_request (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_channel_id INTEGER NOT NULL,
	requester_id VARCHAR(50) NOT NULL,
	advertiser_id VARCHAR(50) NULL,
	backend_message_id VARCHAR(50) NOT NULL UNIQUE,
	message TEXT NOT NULL
	created_at DATETIME NOT NULL,
	resolved_at DATETIME NOT NULL,
	deleted_at DATETIME NULL
);
