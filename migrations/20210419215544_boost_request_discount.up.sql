ALTER TABLE boost_request ADD role_discount_id INTEGER DEFAULT NULL AFTER embed_fields;

ALTER TABLE role_discount DROP INDEX role_discount_unique;

ALTER TABLE role_discount ADD deleted_at DATETIME NULL;
ALTER TABLE role_discount 
	ADD not_deleted TINYINT GENERATED ALWAYS AS (
		IF(deleted_at IS NULL, 1, NULL)
	) VIRTUAL NULL;

CREATE INDEX role_discount_index USING BTREE ON role_discount (guild_id, role_id);
CREATE UNIQUE INDEX role_discount_unique USING BTREE ON role_discount (guild_id, role_id, not_deleted);
