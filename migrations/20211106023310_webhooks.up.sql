CREATE TABLE webhook (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL UNIQUE,
	webhook_url TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	deleted_at DATETIME NULL,
	INDEX(deleted_at)
);

CREATE TABLE webhook_queue (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	webhook_id INTEGER NOT NULL REFERENCES webhook(id),
	body JSON NOT NULL,
	created_at DATETIME NOT NULL,
	INDEX(webhook_id, created_at)
);

CREATE TABLE webhook_attempt (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	webhook_queue_id INTEGER NOT NULL REFERENCES webhook_queue(id),
	status_code INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL,
	INDEX(webhook_queue_id, status_code)
);
