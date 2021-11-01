ALTER TABLE boost_request_channel DROP INDEX boost_request_channel_unique;
CREATE UNIQUE INDEX boost_request_channel_unique ON boost_request_channel (guild_id, frontend_channel_id);

ALTER TABLE boost_request_channel DROP COLUMN is_not_deleted;
ALTER TABLE boost_request_channel DROP COLUMN deleted_at;

