ALTER TABLE boost_request_channel MODIFY COLUMN frontend_channel_id varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL;

DROP TABLE boost_request_preferred_advertiser;

ALTER TABLE boost_request DROP COLUMN advertiser_cut;
ALTER TABLE boost_request DROP COLUMN price;
ALTER TABLE boost_request DROP COLUMN external_id;
