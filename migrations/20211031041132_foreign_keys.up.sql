ALTER TABLE boost_request_role_discount ADD CONSTRAINT brrd_boost_request_id_fk FOREIGN KEY (boost_request_id) REFERENCES boost_request(id);
ALTER TABLE boost_request_role_discount ADD CONSTRAINT brrd_role_discount_id_fk FOREIGN KEY (role_discount_id) REFERENCES role_discount(id);
ALTER TABLE boost_request_preferred_advertiser ADD CONSTRAINT brpa_boost_request_id_fk FOREIGN KEY (boost_request_id) REFERENCES boost_request(id);
