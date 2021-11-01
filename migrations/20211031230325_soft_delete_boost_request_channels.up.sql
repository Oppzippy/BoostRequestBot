ALTER TABLE boost_request_channel ADD deleted_at DATETIME NULL;
ALTER TABLE boost_request_channel ADD is_not_deleted TINYINT(1) GENERATED ALWAYS AS (IF(deleted_at IS NULL, 1, NULL)) VIRTUAL NULL;

ALTER TABLE boost_request_channel DROP INDEX boost_request_channel_unique;
CREATE UNIQUE INDEX boost_request_channel_unique ON boost_request_channel (guild_id, frontend_channel_id, is_not_deleted);
