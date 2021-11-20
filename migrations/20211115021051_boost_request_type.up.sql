ALTER TABLE boost_request ADD request_type VARCHAR(200) NULL AFTER backend_message_id;
ALTER TABLE boost_request ADD discount BIGINT NULL AFTER price;
