CREATE TABLE boost_request_delayed_message (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_id INTEGER NOT NULL REFERENCES boost_request(id),
	delayed_message_id INTEGER NOT NULL REFERENCES delayed_message(id),
	UNIQUE(boost_request_id, delayed_message_id)
);
