ALTER TABLE delayed_message DROP INDEX sent_at;
CREATE INDEX sent_at USING BTREE ON delayed_message (sent_at, send_at);
ALTER TABLE delayed_message DROP COLUMN deleted_at;

DROP TABLE auto_signup_delayed_message;

