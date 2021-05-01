ALTER TABLE boost_request ADD role_discount_id INTEGER DEFAULT NULL AFTER embed_fields;
DROP TABLE role_discount;
CREATE TABLE role_discount (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	guild_id VARCHAR(50) NOT NULL,
	role_id VARCHAR(50) NOT NULL,
	discount DECIMAL(10, 9) NOT NULL,
	created_at DATETIME NOT NULL,
	UNIQUE role_discount_unique (guild_id, role_id)
);

ALTER TABLE role_discount DROP INDEX role_discount_unique;

ALTER TABLE role_discount ADD deleted_at DATETIME NULL;
ALTER TABLE role_discount 
	ADD not_deleted TINYINT GENERATED ALWAYS AS (
		IF(deleted_at IS NULL, 1, NULL)
	) VIRTUAL NULL;

CREATE INDEX role_discount_index USING BTREE ON role_discount (guild_id, role_id);
CREATE UNIQUE INDEX role_discount_unique USING BTREE ON role_discount (guild_id, role_id, not_deleted);
