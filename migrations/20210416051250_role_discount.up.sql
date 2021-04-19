CREATE TABLE role_discount (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	role_id VARCHAR(50) NOT NULL,
	discount DECIMAL(10, 10) NOT NULL,
	created_at DATETIME NOT NULL,
	UNIQUE role_discount_unique (guild_id, role_id)
);
