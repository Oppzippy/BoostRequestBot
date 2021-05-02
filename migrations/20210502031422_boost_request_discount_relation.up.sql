CREATE TABLE boost_request_role_discount (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_id INTEGER NOT NULL,
	role_discount_id INTEGER NOT NULL,
	UNIQUE boost_request_role_discount_unique (boost_request_id, role_discount_id)
);
