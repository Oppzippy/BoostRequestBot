ALTER TABLE boost_request ADD name_visibility ENUM('SHOW', 'SHOW_IN_DMS_ONLY', 'HIDE') DEFAULT 'SHOW' NULL AFTER advertiser_cut;
