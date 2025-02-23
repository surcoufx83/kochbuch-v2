-- +migrate Up
ALTER TABLE
    `users` CHANGE COLUMN `oauth_user_name` `cloudid` VARCHAR(32) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci'
AFTER
    `user_id`,
    CHANGE COLUMN `user_fullname` `clouddisplayname` VARCHAR(128) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci'
AFTER
    `cloudid`,
ADD
    COLUMN `cloudenabled` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0'
AFTER
    `clouddisplayname`,
    CHANGE COLUMN `user_firstname` `firstname` VARCHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci'
AFTER
    `cloudenabled`,
    CHANGE COLUMN `user_lastname` `lastname` VARCHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci'
AFTER
    `firstname`,
    CHANGE COLUMN `user_isactivated` `enabled` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'Wurde das Konto durch einen Admin freigeschaltet'
AFTER
    `lastname`,
    CHANGE COLUMN `user_isadmin` `admin` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0'
AFTER
    `enabled`,
    CHANGE COLUMN `user_email` `email` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci'
AFTER
    `admin`,
    CHANGE COLUMN `user_email_validation` `email_validationphrase` VARCHAR(256) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci'
AFTER
    `email`,
    CHANGE COLUMN `user_email_validated` `email_validated` DATETIME NULL DEFAULT NULL
AFTER
    `email_validationphrase`,
    DROP COLUMN `user_name`,
    DROP COLUMN `user_hash`,
    DROP COLUMN `user_password`,
    DROP COLUMN `user_last_activity`,
    DROP COLUMN `user_registration_completed`,
    DROP COLUMN `user_adconsent`,
    DROP COLUMN `user_betatester`,
    DROP INDEX `user_email`,
ADD
    UNIQUE INDEX `user_email` (`email`) USING BTREE,
    DROP INDEX `user_name`,
ADD
    UNIQUE INDEX `user_name` (`cloudid`) USING BTREE,
    DROP INDEX `oauth_user_name`,
ADD
    UNIQUE INDEX `oauth_user_name` (`cloudid`) USING BTREE;

-- +migrate Down