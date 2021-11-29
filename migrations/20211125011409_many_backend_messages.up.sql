CREATE TABLE boost_request_backend_message (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	boost_request_id INTEGER NOT NULL REFERENCES boost_request(id),
	channel_id VARCHAR(50) NOT NULL,
	message_id VARCHAR(50) NOT NULL,
	INDEX(boost_request_id),
	UNIQUE(channel_id, message_id)
);

INSERT INTO boost_request_backend_message (
	boost_request_id,
	channel_id,
	message_id
) (
	SELECT
		br.id,
		brc.backend_channel_id,
		br.backend_message_id
	FROM
		boost_request AS br
	INNER JOIN boost_request_channel AS brc ON
		br.boost_request_channel_id = brc.id
);

ALTER TABLE boost_request MODIFY COLUMN boost_request_channel_id INTEGER NULL;
ALTER TABLE boost_request MODIFY COLUMN backend_message_id VARCHAR(50) NULL;
ALTER TABLE boost_request ADD guild_id VARCHAR(50) NULL AFTER boost_request_channel_id;

UPDATE
	boost_request AS br
SET
	br.guild_id = (
		SELECT
			brc.guild_id
		FROM
			boost_request_channel AS brc
		WHERE
			brc.id = br.boost_request_channel_id
	);
