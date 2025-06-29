
-- +migrate Up
ALTER TABLE `users`
	CHANGE COLUMN `created` `created` DATETIME NOT NULL DEFAULT current_timestamp() AFTER `email_validated`,
	CHANGE COLUMN `modified` `modified` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() AFTER `created`;

-- +migrate Down
ALTER TABLE `users`
	CHANGE COLUMN `created` `created` DATETIME NULL DEFAULT current_timestamp() AFTER `email_validated`,
	CHANGE COLUMN `modified` `modified` DATETIME NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() AFTER `created`;