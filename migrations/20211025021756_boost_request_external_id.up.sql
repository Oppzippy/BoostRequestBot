ALTER TABLE boost_request ADD external_id VARCHAR(50) NULL AFTER id;
ALTER TABLE boost_request ADD backend_channel_id_override VARCHAR(50) NULL AFTER boost_request_channel_id;
ALTER TABLE boost_request ADD price BIGINT NULL AFTER embed_fields;
ALTER TABLE boost_request ADD advertiser_cut BIGINT NULL AFTER price;

CREATE TABLE boost_request_preferred_advertiser (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_id INTEGER NOT NULL,
	discord_user_id VARCHAR(50) NOT NULL,
	UNIQUE(boost_request_id, discord_user_id)
);
