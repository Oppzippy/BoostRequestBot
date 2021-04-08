CREATE TABLE boost_request_steal_credits (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	user_id VARCHAR(50) NOT NULL,
	credits INTEGER NOT NULL,
	created_at DATETIME NOT NULL,
	INDEX boost_request_steal_credits_index (guild_id, user_id, id)
);
