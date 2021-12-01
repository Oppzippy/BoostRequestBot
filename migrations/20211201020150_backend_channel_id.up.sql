ALTER TABLE boost_request ADD backend_channel_id VARCHAR(50) NULL AFTER guild_id;

UPDATE
	boost_request AS br
SET
	backend_channel_id = (
		SELECT
			brc.backend_channel_id
		FROM
			boost_request_channel AS brc
		WHERE
			brc.id = br.boost_request_channel_id
	);
