ALTER TABLE boost_request DROP COLUMN guild_id;
ALTER TABLE boost_request MODIFY COLUMN backend_message_id VARCHAR(50) NOT NULL;
ALTER TABLE boost_request MODIFY COLUMN boost_request_channel_id INTEGER NOT NULL;

DROP TABLE boost_request_backend_message;
