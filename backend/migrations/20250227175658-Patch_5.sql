-- +migrate Up
CREATE SQL SECURITY INVOKER VIEW `recipes_translationmissing` AS select `recipes`.`recipe_id` AS `recipe_id` from `recipes` where `recipes`.`localized` is null or `recipes`.`localized` < `recipes`.`edited` and `recipes`.`placeholder` = 0;

ALTER TABLE `recipes`
	DROP COLUMN `localized`;

ALTER TABLE `recipes`
	ADD COLUMN `localized` DATETIME NULL DEFAULT NULL AFTER `published`;

UPDATE `recipes` SET `edited` = `created` WHERE `edited` IS NULL;

ALTER SQL SECURITY INVOKER VIEW `allrecipes` AS select `r`.`recipe_id` AS `recipe_id`,`r`.`user_id` AS `user_id`,`r`.`edit_user_id` AS `edit_user_id`,`r`.`aigenerated` AS `aigenerated`,`r`.`placeholder` AS `placeholder`,`r`.`shared_internal` AS `shared_internal`,`r`.`shared_external` AS `shared_external`,`r`.`locale` AS `locale`,`r`.`name_de` AS `name_de`,`r`.`name_en` AS `name_en`,`r`.`name_fr` AS `name_fr`,`r`.`description_de` AS `description_de`,`r`.`description_en` AS `description_en`,`r`.`description_fr` AS `description_fr`,`r`.`servings_count` AS `servings_count`,`r`.`source_description_de` AS `source_description_de`,`r`.`source_description_en` AS `source_description_en`,`r`.`source_description_fr` AS `source_description_fr`,`r`.`source_url` AS `source_url`,`r`.`created` AS `created`,`r`.`edited` AS `edited`,`r`.`modified` AS `modified`,`r`.`localized` AS `localized`,`r`.`published` AS `published`,`r`.`difficulty` AS `difficulty`,`r`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,`p`.`picture_id` AS `picture_id`,`p`.`user_id` AS `picture_user_id`,`p`.`sortindex` AS `picture_sortindex`,`p`.`name_de` AS `picture_name_de`,`p`.`name_en` AS `picture_name_en`,`p`.`name_fr` AS `picture_name_fr`,`p`.`description_de` AS `picture_description_de`,`p`.`description_en` AS `picture_description_en`,`p`.`description_fr` AS `picture_description_fr`,`p`.`filename` AS `picture_filename`,`p`.`fullpath` AS `picture_fullpath`,`p`.`uploaded` AS `picture_uploaded`,`p`.`width` AS `picture_width`,`p`.`height` AS `picture_height`,`rv`.`views` AS `views`,`rc`.`cooked` AS `cooked`,`v`.`votesum` AS `votesum`,`v`.`votes` AS `votes`,`v`.`avgvotes` AS `avgvotes`,`rr`.`votesum` AS `ratesum`,`rr`.`votes` AS `ratings`,`rr`.`avgvotes` AS `avgratings`,`s`.`stepscount` AS `stepscount`,ifnull(`s`.`preparing_time`,-1) AS `preparing_time`,ifnull(`s`.`cooking_time`,-1) AS `cooking_time`,ifnull(`s`.`waiting_time`,-1) AS `waiting_time` from ((((((`recipes` `r` left join `recipe_pictures` `p` on(`p`.`recipe_id` = `r`.`recipe_id` and `p`.`sortindex` = 0)) join `voting_cooked` `rc` on(`rc`.`recipe_id` = `r`.`recipe_id`)) join `voting_views` `rv` on(`rv`.`recipe_id` = `r`.`recipe_id`)) join `voting_hearts` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) join `voting_difficulty` `rr` on(`rr`.`recipe_id` = `r`.`recipe_id`)) join `allrecipessteps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`,`r`.`user_id`,`r`.`edit_user_id`,`r`.`aigenerated`,`r`.`localized`,`r`.`placeholder`,`r`.`shared_internal`,`r`.`shared_external`,`r`.`locale`,`r`.`name_de`,`r`.`name_en`,`r`.`name_fr`,`r`.`description_de`,`r`.`description_en`,`r`.`description_fr`,`r`.`servings_count`,`r`.`source_description_de`,`r`.`source_description_en`,`r`.`source_description_fr`,`r`.`source_url`,`r`.`created`,`r`.`modified`,`r`.`published`,`r`.`difficulty`,`r`.`ingredientsGroupByStep`,`p`.`picture_id`,`p`.`user_id`,`p`.`sortindex`,`p`.`name_de`,`p`.`name_en`,`p`.`name_fr`,`p`.`description_de`,`p`.`description_en`,`p`.`description_fr`,`p`.`filename`,`p`.`fullpath`,`p`.`uploaded`,`p`.`width`,`p`.`height`,`rv`.`views`,`rc`.`cooked`,`v`.`votesum`,`v`.`votes`,`v`.`avgvotes`,`s`.`stepscount`,`s`.`preparing_time`,`s`.`cooking_time`,`s`.`waiting_time`;

-- +migrate Down
DROP VIEW IF EXISTS `recipes_translationmissing`;

ALTER TABLE `recipes`
	DROP COLUMN `localized`;

ALTER TABLE `recipes`
	ADD COLUMN `localized` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `aigenerated`;

