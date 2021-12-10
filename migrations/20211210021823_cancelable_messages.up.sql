CREATE TABLE auto_signup_delayed_message (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	auto_signup_id INTEGER NOT NULL REFERENCES auto_signup_session(id),
	delayed_message_id INTEGER NOT NULL REFERENCES delayed_message(id),
	UNIQUE(auto_signup_id, delayed_message_id)
);

ALTER TABLE delayed_message ADD deleted_at DATETIME NULL AFTER sent_at;
ALTER TABLE delayed_message DROP INDEX sent_at;
CREATE INDEX sent_at USING BTREE ON delayed_message (sent_at, deleted_at, send_at);
