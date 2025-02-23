-- +migrate Up
UPDATE `users` SET `cloudenabled` = 1 WHERE `enabled` = 1;

ALTER TABLE `users`
	DROP INDEX `user_email`;

ALTER TABLE `users`
	CHANGE COLUMN `email` `email` VARCHAR(256) NULL COLLATE 'utf8mb4_general_ci' AFTER `admin`,
	DROP COLUMN `user_avatar`;

UPDATE `users` SET `email` = NULL, `email_validated` = NULL WHERE `email` LIKE 'OAuth2%';

-- +migrate Down
