-- +migrate Up
CREATE SQL SECURITY INVOKER VIEW `recipes_translationmissing` AS select `recipes`.`recipe_id` AS `recipe_id` from `recipes` where `recipes`.`localized` is null or `recipes`.`localized` < `recipes`.`edited`;

ALTER TABLE `recipes`
	DROP COLUMN `localized`;

ALTER TABLE `recipes`
	ADD COLUMN `localized` DATETIME NULL DEFAULT NULL AFTER `published`;

UPDATE `recipes` SET `edited` = `created` WHERE `edited` IS NULL;

-- +migrate Down
DROP VIEW IF EXISTS `recipes_translationmissing`;

ALTER TABLE `recipes`
	DROP COLUMN `localized`;

ALTER TABLE `recipes`
	ADD COLUMN `localized` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `aigenerated`;