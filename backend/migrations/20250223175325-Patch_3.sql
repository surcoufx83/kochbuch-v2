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

UPDATE `recipe_ingredients` SET `ingredient_description_de` = '' WHERE `ingredient_description_de` IS NULL;

UPDATE `recipe_ingredients` SET `ingredient_description_en` = '' WHERE `ingredient_description_en` IS NULL;

ALTER TABLE `recipe_ingredients`
	CHANGE COLUMN `ingredient_quantity` `quantity` DECIMAL(10,3) NULL DEFAULT NULL AFTER `unit_id`,
	CHANGE COLUMN `ingredient_description` `description_de` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `quantity`,
	CHANGE COLUMN `ingredient_description_de` `description_en` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_de`,
	CHANGE COLUMN `ingredient_description_en` `description_fr` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_en`;

ALTER TABLE `recipe_pictures`
	CHANGE COLUMN `picture_sortindex` `sortindex` TINYINT(3) UNSIGNED NOT NULL AFTER `user_id`,
	CHANGE COLUMN `picture_name` `name_de` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `sortindex`,
	ADD COLUMN `name_en` VARCHAR(128) NOT NULL AFTER `name_de`,
	ADD COLUMN `name_fr` VARCHAR(128) NOT NULL AFTER `name_en`,
	CHANGE COLUMN `picture_description` `description_de` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `name_fr`,
	ADD COLUMN `description_en` VARCHAR(1024) NOT NULL AFTER `description_de`,
	ADD COLUMN `description_fr` VARCHAR(1024) NOT NULL AFTER `description_en`,
	CHANGE COLUMN `picture_hash` `hash` VARCHAR(32) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_fr`,
	CHANGE COLUMN `picture_filename` `filename` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `hash`,
	CHANGE COLUMN `picture_full_path` `fullpath` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `filename`,
	CHANGE COLUMN `picture_uploaded` `uploaded` DATETIME NOT NULL DEFAULT current_timestamp() AFTER `fullpath`,
	CHANGE COLUMN `picture_width` `width` SMALLINT UNSIGNED NOT NULL AFTER `uploaded`,
	CHANGE COLUMN `picture_height` `height` SMALLINT UNSIGNED NOT NULL AFTER `width`,
	DROP INDEX `recipe_id_picture_sortindex`,
	ADD UNIQUE INDEX `recipe_id_picture_sortindex` (`recipe_id`, `sortindex`) USING BTREE;

UPDATE `recipe_steps` SET `step_title_de` = '' WHERE `step_title_de` IS NULL;

UPDATE `recipe_steps` SET `step_title_en` = '' WHERE `step_title_en` IS NULL;

UPDATE `recipe_steps` SET `step_data_de` = '' WHERE `step_data_de` IS NULL;

UPDATE `recipe_steps` SET `step_data_en` = '' WHERE `step_data_en` IS NULL;

ALTER TABLE `recipe_steps`
	CHANGE COLUMN `step_no` `sortindex` TINYINT(3) UNSIGNED NOT NULL AFTER `recipe_id`,
	CHANGE COLUMN `step_title` `title_de` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `sortindex`,
	CHANGE COLUMN `step_title_de` `title_en` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `title_de`,
	CHANGE COLUMN `step_title_en` `title_fr` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `title_en`,
	CHANGE COLUMN `step_data` `instruct_de` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `title_fr`,
	CHANGE COLUMN `step_data_de` `instruct_en` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `instruct_de`,
	CHANGE COLUMN `step_data_en` `instruct_fr` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `instruct_en`,
	CHANGE COLUMN `step_time_preparation` `preparing` SMALLINT NULL DEFAULT NULL AFTER `instruct_fr`,
	CHANGE COLUMN `step_time_cooking` `cooking` SMALLINT NULL DEFAULT NULL AFTER `preparing`,
	CHANGE COLUMN `step_time_chill` `waiting` SMALLINT NULL DEFAULT NULL AFTER `cooking`,
	DROP INDEX `recipe_id_step_no`,
	ADD UNIQUE INDEX `recipe_id_step_no` (`recipe_id`, `sortindex`) USING BTREE;

ALTER SQL SECURITY INVOKER VIEW `allrecipessteps` AS select `r`.`recipe_id` AS `recipe_id`,count(`s`.`sortindex`) AS `stepscount`,sum(if(ifnull(`s`.`preparing`,0) < 0,0,`s`.`preparing`)) AS `preparing_time`,sum(if(ifnull(`s`.`cooking`,0) < 0,0,`s`.`cooking`)) AS `cooking_time`,sum(if(ifnull(`s`.`waiting`,0) < 0,0,`s`.`waiting`)) AS `waiting_time` from (`recipes` `r` left join `recipe_steps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id` ;

ALTER SQL SECURITY INVOKER VIEW `allrecipes` AS select `r`.`recipe_id` AS `recipe_id`,`r`.`user_id` AS `user_id`,`r`.`edit_user_id` AS `edit_user_id`,`r`.`aigenerated` AS `aigenerated`,`r`.`localized` AS `localized`,`r`.`placeholder` AS `placeholder`,`r`.`shared_internal` AS `shared_internal`,`r`.`shared_external` AS `shared_external`,`r`.`locale` AS `locale`,`r`.`name_de` AS `name_de`,`r`.`name_en` AS `name_en`,`r`.`name_fr` AS `name_fr`,`r`.`description_de` AS `description_de`,`r`.`description_en` AS `description_en`,`r`.`description_fr` AS `description_fr`,`r`.`servings_count` AS `servings_count`,`r`.`source_description_de` AS `source_description_de`,`r`.`source_description_en` AS `source_description_en`,`r`.`source_description_fr` AS `source_description_fr`,`r`.`source_url` AS `source_url`,`r`.`created` AS `created`,`r`.`modified` AS `modified`,`r`.`published` AS `published`,`r`.`difficulty` AS `difficulty`,`r`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,`p`.`picture_id` AS `picture_id`,`p`.`user_id` AS `picture_user_id`,`p`.`sortindex` AS `picture_sortindex`,`p`.`name_de` AS `picture_name_de`,`p`.`name_en` AS `picture_name_en`,`p`.`name_fr` AS `picture_name_fr`,`p`.`description_de` AS `picture_description_de`,`p`.`description_en` AS `picture_description_en`,`p`.`description_fr` AS `picture_description_fr`,`p`.`filename` AS `picture_filename`,`p`.`fullpath` AS `picture_fullpath`,`p`.`uploaded` AS `picture_uploaded`,`p`.`width` AS `picture_width`,`p`.`height` AS `picture_height`,`rv`.`views` AS `views`,`rc`.`cooked` AS `cooked`,`v`.`votesum` AS `votesum`,`v`.`votes` AS `votes`,`v`.`avgvotes` AS `avgvotes`,`rr`.`votesum` AS `ratesum`,`rr`.`votes` AS `ratings`,`rr`.`avgvotes` AS `avgratings`,`s`.`stepscount` AS `stepscount`,ifnull(`s`.`preparing_time`,-1) AS `preparing_time`,ifnull(`s`.`cooking_time`,-1) AS `cooking_time`,ifnull(`s`.`waiting_time`,-1) AS `waiting_time` from ((((((`recipes` `r` left join `recipe_pictures` `p` on(`p`.`recipe_id` = `r`.`recipe_id` and `p`.`sortindex` = 0)) join `voting_cooked` `rc` on(`rc`.`recipe_id` = `r`.`recipe_id`)) join `voting_views` `rv` on(`rv`.`recipe_id` = `r`.`recipe_id`)) join `voting_hearts` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) join `voting_difficulty` `rr` on(`rr`.`recipe_id` = `r`.`recipe_id`)) join `allrecipessteps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`,`r`.`user_id`,`r`.`edit_user_id`,`r`.`aigenerated`,`r`.`localized`,`r`.`placeholder`,`r`.`shared_internal`,`r`.`shared_external`,`r`.`locale`,`r`.`name_de`,`r`.`name_en`,`r`.`name_fr`,`r`.`description_de`,`r`.`description_en`,`r`.`description_fr`,`r`.`servings_count`,`r`.`source_description_de`,`r`.`source_description_en`,`r`.`source_description_fr`,`r`.`source_url`,`r`.`created`,`r`.`modified`,`r`.`published`,`r`.`difficulty`,`r`.`ingredientsGroupByStep`,`p`.`picture_id`,`p`.`user_id`,`p`.`sortindex`,`p`.`name_de`,`p`.`name_en`,`p`.`name_fr`,`p`.`description_de`,`p`.`description_en`,`p`.`description_fr`,`p`.`filename`,`p`.`fullpath`,`p`.`uploaded`,`p`.`width`,`p`.`height`,`rv`.`views`,`rc`.`cooked`,`v`.`votesum`,`v`.`votes`,`v`.`avgvotes`,`s`.`stepscount`,`s`.`preparing_time`,`s`.`cooking_time`,`s`.`waiting_time` ;

DROP VIEW `allrecipes_nouser`;

DROP VIEW `recipes_my`;

DROP VIEW `apicache_view`;

DROP TABLE `apicache`;

-- +migrate Down
ALTER TABLE `users`
	CHANGE COLUMN `cloudid` `cloudid` VARCHAR(32) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci' AFTER `user_id`,
    DROP COLUMN `cloudsync`,
    DROP COLUMN `cloudsync_status`,
	DROP COLUMN `created`,
	DROP COLUMN `modified`;

DROP TABLE `user_groups`;

DROP TABLE `groups`;

ALTER TABLE `recipe_ingredients`
	CHANGE COLUMN `quantity` `ingredient_quantity` DECIMAL(10,3) NULL DEFAULT NULL AFTER `unit_id`,
	CHANGE COLUMN `description_de` `ingredient_description` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `ingredient_quantity`,
	CHANGE COLUMN `description_en` `ingredient_description_de` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `ingredient_description`,
	CHANGE COLUMN `description_fr` `ingredient_description_en` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `ingredient_description_de`;

ALTER TABLE `recipe_pictures`
	CHANGE COLUMN `sortindex` `picture_sortindex` TINYINT(3) UNSIGNED NOT NULL AFTER `user_id`,
	CHANGE COLUMN `name_de` `picture_name` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `picture_sortindex`,
	DROP COLUMN `name_en`,
	DROP COLUMN `name_fr`,
	CHANGE COLUMN `description_de` `picture_description` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `picture_name`,
	DROP COLUMN `description_en`,
	DROP COLUMN `description_fr`,
	CHANGE COLUMN `hash` `picture_hash` VARCHAR(32) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `picture_description`,
	CHANGE COLUMN `filename` `picture_filename` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `picture_hash`,
	CHANGE COLUMN `fullpath` `picture_full_path` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `picture_filename`,
	CHANGE COLUMN `uploaded` `picture_uploaded` DATETIME NOT NULL DEFAULT current_timestamp() AFTER `picture_full_path`,
	CHANGE COLUMN `width` `picture_width` SMALLINT UNSIGNED NOT NULL AFTER `picture_uploaded`,
	CHANGE COLUMN `height` `picture_height` SMALLINT UNSIGNED NOT NULL AFTER `picture_width`,
	DROP INDEX `recipe_id_picture_sortindex`,
	ADD UNIQUE INDEX `recipe_id_picture_sortindex` (`recipe_id`, `picture_sortindex`) USING BTREE;

ALTER TABLE `recipe_steps`
	CHANGE COLUMN `sortindex` `step_no` TINYINT(3) UNSIGNED NOT NULL AFTER `recipe_id`,
	CHANGE COLUMN `title_de` `step_title` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `step_no`,
	CHANGE COLUMN `title_en` `step_title_de` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `step_title`,
	CHANGE COLUMN `title_fr` `step_title_en` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `step_title_de`,
	CHANGE COLUMN `instruct_de` `step_data` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `step_title_en`,
	CHANGE COLUMN `instruct_en` `step_data_de` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `step_data`,
	CHANGE COLUMN `instruct_fr` `step_data_en` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `step_data_de`,
	CHANGE COLUMN `preparing` `step_time_preparation` SMALLINT NULL DEFAULT NULL AFTER `step_data_en`,
	CHANGE COLUMN `cooking` `step_time_cooking` SMALLINT NULL DEFAULT NULL AFTER `step_time_preparation`,
	CHANGE COLUMN `waiting` `step_time_chill` SMALLINT NULL DEFAULT NULL AFTER `step_time_cooking`,
	DROP INDEX `recipe_id_step_no`,
	ADD UNIQUE INDEX `recipe_id_step_no` (`recipe_id`, `step_no`) USING BTREE;

ALTER SQL SECURITY INVOKER VIEW `allrecipessteps` AS select `r`.`recipe_id` AS `recipe_id`,count(`s`.`step_no`) AS `stepscount`,sum(if(ifnull(`s`.`step_time_preparation`,0) < 0,0,`s`.`step_time_preparation`)) AS `preparationtime`,sum(if(ifnull(`s`.`step_time_cooking`,0) < 0,0,`s`.`step_time_cooking`)) AS `cookingtime`,sum(if(ifnull(`s`.`step_time_chill`,0) < 0,0,`s`.`step_time_chill`)) AS `chilltime` from (`recipes` `r` left join `recipe_steps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id` ;

ALTER SQL SECURITY INVOKER VIEW `allrecipes` AS select `r`.`recipe_id` AS `recipe_id`,ifnull(`r`.`user_id`,0) AS `user_id`,ifnull(`r`.`edit_user_id`,0) AS `edit_user_id`,`r`.`aigenerated` AS `aigenerated`,`r`.`localized` AS `localized`,`r`.`placeholder` AS `placeholder`,`r`.`shared_internal` AS `shared_internal`,`r`.`shared_external` AS `shared_external`,`r`.`locale` AS `locale`,`r`.`name_de` AS `name_de`,`r`.`name_en` AS `name_en`,`r`.`name_fr` AS `name_fr`,`r`.`description_de` AS `description_de`,`r`.`description_en` AS `description_en`,`r`.`description_fr` AS `description_fr`,`r`.`servings_count` AS `servings_count`,`r`.`source_description_de` AS `source_description_de`,`r`.`source_description_en` AS `source_description_en`,`r`.`source_description_fr` AS `source_description_fr`,`r`.`source_url` AS `source_url`,`r`.`created` AS `created`,`r`.`modified` AS `modified`,`r`.`published` AS `published`,`r`.`difficulty` AS `difficulty`,`r`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,ifnull(`p`.`picture_id`,0) AS `picture_id`,ifnull(`p`.`picture_sortindex`,0) AS `picture_sortindex`,ifnull(`p`.`picture_name`,'') AS `picture_name`,ifnull(`p`.`picture_description`,'') AS `picture_description`,ifnull(`p`.`picture_hash`,'') AS `picture_hash`,ifnull(`p`.`picture_filename`,'') AS `picture_filename`,ifnull(`p`.`picture_full_path`,'') AS `picture_full_path`,cast(ifnull(`p`.`picture_uploaded`,'2000-01-01 00:00:00') as datetime) AS `picture_uploaded`,ifnull(`p`.`picture_width`,0) AS `picture_width`,ifnull(`p`.`picture_height`,0) AS `picture_height`,`rv`.`views` AS `views`,`rc`.`cooked` AS `cooked`,`v`.`votesum` AS `votesum`,`v`.`votes` AS `votes`,`v`.`avgvotes` AS `avgvotes`,`rr`.`votesum` AS `ratesum`,`rr`.`votes` AS `ratings`,`rr`.`avgvotes` AS `avgratings`,`s`.`stepscount` AS `stepscount`,ifnull(`s`.`preparationtime`,-1) AS `preparationtime`,ifnull(`s`.`cookingtime`,-1) AS `cookingtime`,ifnull(`s`.`chilltime`,-1) AS `chilltime` from ((((((`recipes` `r` left join `recipe_pictures` `p` on(`p`.`recipe_id` = `r`.`recipe_id` and `p`.`picture_sortindex` = 0)) join `voting_cooked` `rc` on(`rc`.`recipe_id` = `r`.`recipe_id`)) join `voting_views` `rv` on(`rv`.`recipe_id` = `r`.`recipe_id`)) join `voting_hearts` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) join `voting_difficulty` `rr` on(`rr`.`recipe_id` = `r`.`recipe_id`)) join `allrecipessteps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`,`r`.`user_id`,`r`.`edit_user_id`,`r`.`aigenerated`,`r`.`localized`,`r`.`placeholder`,`r`.`shared_internal`,`r`.`shared_external`,`r`.`locale`,`r`.`name_de`,`r`.`name_en`,`r`.`name_fr`,`r`.`description_de`,`r`.`description_en`,`r`.`description_fr`,`r`.`servings_count`,`r`.`source_description_de`,`r`.`source_description_en`,`r`.`source_description_fr`,`r`.`source_url`,`r`.`created`,`r`.`modified`,`r`.`published`,`r`.`difficulty`,`r`.`ingredientsGroupByStep`,`p`.`picture_id`,`p`.`picture_sortindex`,`p`.`picture_name`,`p`.`picture_description`,`p`.`picture_hash`,`p`.`picture_filename`,`p`.`picture_full_path`,`p`.`picture_uploaded`,`p`.`picture_width`,`p`.`picture_height`,`rv`.`views`,`rc`.`cooked`,`v`.`votesum`,`v`.`votes`,`v`.`avgvotes`,`s`.`stepscount`,`s`.`preparationtime`,`s`.`cookingtime`,`s`.`chilltime` ;

CREATE SQL SECURITY INVOKER VIEW `allrecipes_nouser` AS select * from `allrecipes` where `allrecipes`.`shared_external` = 1 ;

CREATE SQL SECURITY INVOKER VIEW `recipes_my` AS select `recipes`.`recipe_id` AS `recipe_id`,`recipes`.`user_id` AS `user_id`,`kochbuch`.`recipes`.`name_de` AS `recipe_name`,`kochbuch`.`recipes`.`shared_internal` AS `recipe_public_internal`,`kochbuch`.`recipes`.`shared_external` AS `recipe_public_external` from `recipes` ;

CREATE TABLE `apicache` (
	`route` VARCHAR(512) NOT NULL COLLATE 'utf8mb4_general_ci',
 	`ts` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
 	`response` TEXT NOT NULL COLLATE 'utf8mb4_general_ci',
 	PRIMARY KEY (`route`, `ts`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB;

CREATE SQL SECURITY INVOKER VIEW `apicache_view` AS select `a`.`route` AS `route`,`a`.`ts` AS `ts`,`a`.`response` AS `response` from (`apicache` `a` join `lastactivityview` `l` on(`l`.`ts` = `a`.`ts`)) ;
