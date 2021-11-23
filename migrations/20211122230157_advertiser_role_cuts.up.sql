CREATE TABLE boost_request_role_cut (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_id INTEGER NOT NULL REFERENCES boost_request(id),
	role_id VARCHAR(50) NOT NULL,
	role_cut BIGINT NOT NULL,
	UNIQUE(boost_request_id, role_id)
);
