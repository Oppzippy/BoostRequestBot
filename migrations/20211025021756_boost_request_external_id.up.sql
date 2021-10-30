ALTER TABLE boost_request ADD external_id VARCHAR(50) NULL AFTER id;
ALTER TABLE boost_request ADD price BIGINT NULL AFTER embed_fields;
ALTER TABLE boost_request ADD advertiser_cut BIGINT NULL AFTER price;

CREATE TABLE boost_request_preferred_advertiser (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_id INTEGER NOT NULL,
	discord_user_id VARCHAR(50) NOT NULL,
	UNIQUE(boost_request_id, discord_user_id)
);

ALTER TABLE boost_request_channel MODIFY COLUMN frontend_channel_id varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL;
