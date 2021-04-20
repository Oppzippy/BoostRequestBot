ALTER TABLE boost_request DROP COLUMN role_discount_id;

ALTER TABLE role_discount DROP INDEX role_discount_index;
ALTER TABLE role_discount DROP INDEX role_discount_unique;
ALTER TABLE role_discount DROP COLUMN not_deleted;
ALTER TABLE role_discount DROP COLUMN deleted_at;

CREATE UNIQUE INDEX role_discount_unique USING BTREE ON role_discount (guild_id, role_id);
