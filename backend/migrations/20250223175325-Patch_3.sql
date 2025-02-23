-- +migrate Up
ALTER TABLE `users`
    CHANGE COLUMN `cloudid` `cloudid` VARCHAR(32) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `user_id`,
    ADD COLUMN `cloudsync` DATETIME NULL DEFAULT NULL AFTER `cloudenabled`,
    ADD COLUMN `cloudsync_status` SMALLINT UNSIGNED NOT NULL DEFAULT 0 AFTER `cloudsync`,
	ADD COLUMN `created` DATETIME NULL DEFAULT current_timestamp() AFTER `email_validated`,
	ADD COLUMN `modified` DATETIME NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() AFTER `created`;

CREATE TABLE `groups` (
	`id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`displayname` VARCHAR(64) NOT NULL,
	`ncname` VARCHAR(64) NOT NULL,
	`access_granted` TINYINT UNSIGNED NOT NULL DEFAULT 0,
	`is_admin` TINYINT UNSIGNED NOT NULL DEFAULT 0,
	`created` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	PRIMARY KEY (`id`),
	UNIQUE INDEX `ncname` (`ncname`)
)
COLLATE='utf8mb4_general_ci';

CREATE TABLE `user_groups` (
	`id` MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`userid` MEDIUMINT UNSIGNED NOT NULL,
	`groupid` SMALLINT UNSIGNED NOT NULL,
	`granted` DATETIME NOT NULL DEFAULT current_timestamp(),
	PRIMARY KEY (`id`),
	UNIQUE INDEX `userid_groupid` (`userid`, `groupid`),
	CONSTRAINT `FK__users` FOREIGN KEY (`userid`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT `FK__groups` FOREIGN KEY (`groupid`) REFERENCES `groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci';

-- +migrate Down
ALTER TABLE `users`
    CHANGE COLUMN `cloudid` `cloudid` VARCHAR(32) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci' AFTER `user_id`,
    DROP COLUMN `cloudsync`,
    DROP COLUMN `cloudsync_status`,
	DROP COLUMN `created`,
	DROP COLUMN `modified`;

DROP TABLE `groups`;

DROP TABLE `user_groups`;