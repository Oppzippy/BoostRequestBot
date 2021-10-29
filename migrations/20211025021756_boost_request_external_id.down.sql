DROP TABLE boost_request_preferred_advertiser;

ALTER TABLE boost_request DROP COLUMN advertiser_cut;
ALTER TABLE boost_request DROP COLUMN price;
ALTER TABLE boost_request DROP COLUMN backend_channel_id_override;
ALTER TABLE boost_request DROP COLUMN external_id;
