CREATE TABLE delayed_message (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	destination_id VARCHAR(50) NOT NULL,
	destination_type ENUM('CHANNEL', 'USER') NOT NULL,
	fallback_channel_id VARCHAR(50),
	message_json JSON NOT NULL,
	send_at DATETIME NOT NULL,
	sent_at DATETIME,
	INDEX(sent_at, send_at)
);
