ALTER TABLE boost_request DROP COLUMN role_discount_id;
DROP TABLE role_discount;

CREATE TABLE role_discount (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	role_id VARCHAR(50) NOT NULL,
	boost_type VARCHAR(200) NOT NULL,
	discount DECIMAL(10, 9),
	created_at DATETIME NOT NULL,
	deleted_at DATETIME NULL,
	not_deleted TINYINT GENERATED ALWAYS AS (
		IF(deleted_at IS NULL, 1, NULL)
	) VIRTUAL NULL,
	UNIQUE role_discount_unique (guild_id, role_id, boost_type, not_deleted)
);
