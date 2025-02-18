# Database Changes to Version 1

## Updated Units

### Modify units table

```sql
UPDATE `units` SET `unit_name_de` = '' WHERE `unit_name_de` IS NULL;

UPDATE `units` SET `unit_name_en` = '' WHERE `unit_name_en` IS NULL;

ALTER TABLE `units`
    CHANGE COLUMN `unit_name` `sg_name_de` VARCHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `localized`,
    CHANGE COLUMN `unit_name_de` `sg_name_en` VARCHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `sg_name_de`,
    CHANGE COLUMN `unit_name_en` `sg_name_fr` VARCHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `sg_name_en`,
    ADD COLUMN `pl_name_de` VARCHAR(64) NOT NULL DEFAULT '' AFTER `sg_name_fr`,
    ADD COLUMN `pl_name_en` VARCHAR(64) NOT NULL DEFAULT '' AFTER `pl_name_de`,
    ADD COLUMN `pl_name_fr` VARCHAR(64) NOT NULL DEFAULT '' AFTER `pl_name_en`,
    ADD COLUMN `decimal_places` TINYINT UNSIGNED NOT NULL DEFAULT 0 AFTER `pl_name_fr`,
    ADD COLUMN `fractional` TINYINT UNSIGNED NOT NULL DEFAULT 0 AFTER `decimal_places`,
    DROP INDEX `unit_name`;

UPDATE `units` SET `pl_name_de` = `sg_name_de` WHERE `pl_name_de` = '';

ALTER TABLE `units`
    DROP FOREIGN KEY `FK_units_units`;

ALTER TABLE `units`
    CHANGE COLUMN `useunit_id` `supersededby_unitid` SMALLINT(5) UNSIGNED NULL DEFAULT NULL AFTER `unit_id`,
    ADD COLUMN `saveas_unitid` SMALLINT(5) UNSIGNED NULL DEFAULT NULL AFTER `supersededby_unitid`,
    ADD COLUMN `saveas_factor` DECIMAL(20,6) UNSIGNED NOT NULL DEFAULT 1 AFTER `saveas_unitid`,
    DROP INDEX `FK_units_units`,
    ADD INDEX `FK_units_units` (`supersededby_unitid`) USING BTREE,
    ADD CONSTRAINT `FK_units_units` FOREIGN KEY (`supersededby_unitid`) REFERENCES `units` (`unit_id`) ON UPDATE CASCADE ON DELETE SET NULL,
    ADD CONSTRAINT `FK_units_units_2` FOREIGN KEY (`saveas_unitid`) REFERENCES `units` (`unit_id`) ON UPDATE CASCADE ON DELETE SET NULL;

UPDATE `units` SET `sg_name_en` = `sg_name_de`, `sg_name_fr` = `sg_name_de`, `pl_name_en` = `sg_name_de`, `pl_name_fr` = `sg_name_de` WHERE `unit_id` = 1 AND `sg_name_de` = 'g';
UPDATE `units` SET `saveas_unitid` = 1, `saveas_factor` = 1000.0, `sg_name_en` = `sg_name_de`, `sg_name_fr` = `sg_name_de`, `pl_name_en` = `sg_name_de`, `pl_name_fr` = `sg_name_de`, `decimal_places` = 2 WHERE `unit_id` = 2 AND `sg_name_de` = 'kg';
UPDATE `units` SET `saveas_unitid` = 6, `saveas_factor` = 1000.0, `sg_name_en` = `sg_name_de`, `sg_name_fr` = `sg_name_de`, `pl_name_en` = `sg_name_de`, `pl_name_fr` = `sg_name_de`, `decimal_places` = 2 WHERE `unit_id` = 3 AND `sg_name_de` = 'l';
UPDATE `units` SET `sg_name_en` = 'pc', `sg_name_fr` = 'pce', `pl_name_en` = 'pcs', `pl_name_fr` = 'pces', `fractional` = 1 WHERE `unit_id` = 4 AND `sg_name_de` = 'St.';
UPDATE `units` SET `sg_name_de` = 'Zweig', `sg_name_en` = 'sprig', `sg_name_fr` = 'brin', `pl_name_de` = 'Zweige', `pl_name_en` = 'sprigs', `pl_name_fr` = 'brins' WHERE `unit_id` = 5 AND `sg_name_de` = 'Zweige';
UPDATE `units` SET `sg_name_en` = `sg_name_de`, `sg_name_fr` = `sg_name_de`, `pl_name_en` = `sg_name_de`, `pl_name_fr` = `sg_name_de` WHERE `unit_id` = 6 AND `sg_name_de` = 'ml';
UPDATE `units` SET `sg_name_de` = 'EL', `sg_name_en` = 'tbsp', `sg_name_fr` = 'c. à s.', `pl_name_de` = 'EL', `pl_name_en` = 'tbsp', `pl_name_fr` = 'c. à s.' WHERE `unit_id` = 7 AND `sg_name_de` = 'EL';
UPDATE `units` SET `sg_name_de` = 'TL', `sg_name_en` = 'tsp', `sg_name_fr` = 'c. à c.', `pl_name_de` = 'TL', `pl_name_en` = 'tsp', `pl_name_fr` = 'c. à c.' WHERE `unit_id` = 8 AND `sg_name_de` = 'TL';
UPDATE `units` SET `sg_name_de` = 'Msp.', `sg_name_en` = 'pinch', `sg_name_fr` = 'pincée', `pl_name_de` = 'Msp.', `pl_name_en` = 'pinches', `pl_name_fr` = 'pincées' WHERE `unit_id` = 9 AND `sg_name_de` = 'Msp.';
UPDATE `units` SET `sg_name_de` = 'Blatt', `sg_name_en` = 'leaf', `sg_name_fr` = 'feuille', `pl_name_de` = 'Blätter', `pl_name_en` = 'leaves', `pl_name_fr` = 'feuilles' WHERE `unit_id` = 52 AND `sg_name_de` = 'Blätter';
UPDATE `units` SET `sg_name_de` = 'Zehe', `sg_name_en` = 'clove', `sg_name_fr` = 'gousse', `pl_name_de` = 'Zehen', `pl_name_en` = 'cloves', `pl_name_fr` = 'gousses' WHERE `unit_id` = 53 AND `sg_name_de` = 'Zehen';
UPDATE `units` SET `sg_name_de` = 'Becher', `sg_name_en` = 'cup', `sg_name_fr` = 'gobelet', `pl_name_de` = 'Becher', `pl_name_en` = 'cups', `pl_name_fr` = 'gobelets', `fractional` = 1 WHERE `unit_id` = 54 AND `sg_name_de` = 'Becher';
UPDATE `units` SET `sg_name_de` = 'Tüte', `sg_name_en` = 'bag', `sg_name_fr` = 'sachet', `pl_name_de` = 'Tüten', `pl_name_en` = 'bags', `pl_name_fr` = 'sachets', `fractional` = 1 WHERE `unit_id` = 55 AND `sg_name_de` = 'Tüten';
UPDATE `units` SET `sg_name_de` = 'Scheibe', `sg_name_en` = 'slice', `sg_name_fr` = 'tranche', `pl_name_de` = 'Scheiben', `pl_name_en` = 'slices', `pl_name_fr` = 'tranches' WHERE `unit_id` = 56 AND `sg_name_de` = 'Scheiben';
UPDATE `units` SET `sg_name_de` = 'Dose', `pl_name_de` = 'Dosen', `sg_name_en` = 'can', `pl_name_en` = 'cans', `sg_name_fr` = 'boîte', `pl_name_fr` = 'boîtes', `fractional` = 1 WHERE `unit_id` = 57 AND `sg_name_de` = 'Dose';
UPDATE `units` SET `sg_name_de` = 'Priese', `pl_name_de` = 'Priesen', `sg_name_en` = 'pinch', `pl_name_en` = 'pinches', `sg_name_fr` = 'pincée', `pl_name_fr` = 'pincées' WHERE `unit_id` = 58 AND `sg_name_de` = 'Priese';
UPDATE `units` SET `sg_name_de` = 'nach Belieben', `pl_name_de` = 'nach Belieben', `sg_name_en` = 'to taste', `pl_name_en` = 'to taste', `sg_name_fr` = 'selon goût', `pl_name_fr` = 'selon goût' WHERE `unit_id` = 59 AND `sg_name_de` = 'nach Belieben';
UPDATE `units` SET `sg_name_de` = 'Würfel', `pl_name_de` = 'Würfel', `sg_name_en` = 'cube', `pl_name_en` = 'cubes', `sg_name_fr` = 'cube', `pl_name_fr` = 'cubes' WHERE `unit_id` = 60 AND `sg_name_de` = 'Würfel';
UPDATE `units` SET `sg_name_de` = 'Packung', `pl_name_de` = 'Packungen', `sg_name_en` = 'pack', `pl_name_en` = 'packs', `sg_name_fr` = 'paquet', `pl_name_fr` = 'paquets', `fractional` = 1 WHERE `unit_id` = 61 AND `sg_name_de` = 'Packung';
UPDATE `units` SET `sg_name_de` = 'Bund', `pl_name_de` = 'Bunde', `sg_name_en` = 'brunch', `pl_name_en` = 'brunches', `sg_name_fr` = 'botte', `pl_name_fr` = 'bottes', `fractional` = 1 WHERE `unit_id` = 63 AND `sg_name_de` = 'Bund';



UPDATE `recipe_ingredients` SET `unit_id` = 4 WHERE `unit_id` = 50; -- St -> St.

UPDATE `recipe_ingredients` SET `unit_id` = 5 WHERE `unit_id` = 10; -- Zweig -> Zweige
DELETE FROM `units` WHERE `unit_id` = 10 AND `sg_name_de` = 'Zweig';

UPDATE `recipe_ingredients` SET `unit_id` = 53 WHERE `unit_id` = 62; -- Zehe -> Zehen
DELETE FROM `units` WHERE `unit_id` = 62 AND `sg_name_de` = 'Zehe';

```

### Create a view unitsview for better access

```sql
ALGORITHM = UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `unitsview` AS select `units`.`unit_id` AS `unit_id`,ifnull(`units`.`supersededby_unitid`,0) AS `supersededby_unitid`,ifnull(`units`.`saveas_unitid`,0) AS `saveas_unitid`,`units`.`saveas_factor` AS `saveas_factor`,`units`.`localized` AS `localized`,`units`.`sg_name_de` AS `sg_name_de`,`units`.`sg_name_en` AS `sg_name_en`,`units`.`sg_name_fr` AS `sg_name_fr`,`units`.`pl_name_de` AS `pl_name_de`,`units`.`pl_name_en` AS `pl_name_en`,`units`.`pl_name_fr` AS `pl_name_fr`,`units`.`decimal_places` AS `decimal_places`,`units`.`fractional` AS `fractional`,`units`.`created` AS `created`,`units`.`updated` AS `updated` from `units`;
```

## Updated Recipes

### Modify recipes table

```sql
UPDATE `recipes` SET `recipe_name_de` = '' WHERE `recipe_name_de` IS NULL;
UPDATE `recipes` SET `recipe_name_en` = '' WHERE `recipe_name_en` IS NULL;
UPDATE `recipes` SET `recipe_description_de` = '' WHERE `recipe_description_de` IS NULL;
UPDATE `recipes` SET `recipe_description_en` = '' WHERE `recipe_description_en` IS NULL;

ALTER TABLE `recipe_pictures`
    CHANGE COLUMN `picture_uploaded` `picture_uploaded` DATETIME NOT NULL DEFAULT current_timestamp() AFTER `picture_full_path`;

ALTER TABLE `recipes`
    CHANGE COLUMN `aigenerated` `aigenerated` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `edit_user_id`,
    CHANGE COLUMN `localized` `localized` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `aigenerated`,
    CHANGE COLUMN `recipe_placeholder` `placeholder` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `localized`,
    CHANGE COLUMN `recipe_public_internal` `shared_internal` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `placeholder`,
    CHANGE COLUMN `recipe_public_external` `shared_external` TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER `shared_internal`,
    CHANGE COLUMN `recipe_name` `name_de` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `shared_external`,
    CHANGE COLUMN `recipe_name_de` `name_en` VARCHAR(256) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci' AFTER `name_de`,
    CHANGE COLUMN `recipe_name_en` `name_fr` VARCHAR(256) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci' AFTER `name_en`,
    CHANGE COLUMN `recipe_description` `description_de` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `name_fr`,
    CHANGE COLUMN `recipe_description_de` `description_en` VARCHAR(1024) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_de`,
    CHANGE COLUMN `recipe_description_en` `description_fr` VARCHAR(1024) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_en`,
    CHANGE COLUMN `recipe_eater` `servings_count` TINYINT(3) UNSIGNED NOT NULL AFTER `description_fr`,
    CHANGE COLUMN `recipe_source_desc` `source_description_de` VARCHAR(1024) NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `servings_count`,
    ADD COLUMN `source_description_en` VARCHAR(1024) NULL DEFAULT '' AFTER `source_description_de`,
    ADD COLUMN `source_description_fr` VARCHAR(1024) NULL DEFAULT '' AFTER `source_description_en`,
    CHANGE COLUMN `recipe_source_url` `source_url` VARCHAR(256) NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `source_description_fr`,
    DROP INDEX `FK_recipes_users`,
    DROP FOREIGN KEY `FK_recipes_users`;

ALTER TABLE `recipes`
    CHANGE COLUMN `recipe_created` `created` TIMESTAMP NOT NULL DEFAULT current_timestamp() AFTER `source_url`,
    CHANGE COLUMN `recipe_edited` `edited` TIMESTAMP NULL DEFAULT NULL AFTER `created`,
    CHANGE COLUMN `recipe_modified` `modified` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() AFTER `edited`,
    CHANGE COLUMN `recipe_published` `published` TIMESTAMP NULL DEFAULT NULL AFTER `modified`,
    CHANGE COLUMN `recipe_difficulty` `difficulty` ENUM('0','1','2','3') NOT NULL DEFAULT '0' COLLATE 'utf8mb4_general_ci' AFTER `published`,
    ADD CONSTRAINT `FK_recipes_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE ON DELETE SET NULL,
    ADD CONSTRAINT `FK_recipes_users_2` FOREIGN KEY (`edit_user_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE `recipes`
    CHANGE COLUMN `name_en` `name_en` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `name_de`,
    CHANGE COLUMN `name_fr` `name_fr` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `name_en`,
    CHANGE COLUMN `description_en` `description_en` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_de`,
    CHANGE COLUMN `description_fr` `description_fr` VARCHAR(1024) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `description_en`,
    CHANGE COLUMN `source_description_de` `source_description_de` VARCHAR(1024) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `servings_count`,
    CHANGE COLUMN `source_description_en` `source_description_en` VARCHAR(1024) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `source_description_de`,
    CHANGE COLUMN `source_description_fr` `source_description_fr` VARCHAR(1024) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `source_description_en`,
    CHANGE COLUMN `source_url` `source_url` VARCHAR(256) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `source_description_fr`;

ALTER TABLE `recipes`
    CHANGE COLUMN `created` `created` DATETIME NOT NULL DEFAULT current_timestamp() AFTER `source_url`,
    CHANGE COLUMN `edited` `edited` DATETIME NULL DEFAULT NULL AFTER `created`,
    CHANGE COLUMN `modified` `modified` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() AFTER `edited`,
    CHANGE COLUMN `published` `published` DATETIME NULL DEFAULT NULL AFTER `modified`;

ALTER TABLE `recipes`
    ADD COLUMN `locale` VARCHAR(2) NOT NULL DEFAULT 'de' AFTER `shared_external`;

ALTER ALGORITHM = UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `allrecipes` AS select `r`.`recipe_id` AS `recipe_id`,ifnull(`r`.`user_id`,0) AS `user_id`,ifnull(`r`.`edit_user_id`,0) AS `edit_user_id`,`r`.`aigenerated` AS `aigenerated`,`r`.`localized` AS `localized`,`r`.`placeholder` AS `placeholder`,`r`.`shared_internal` AS `shared_internal`,`r`.`shared_external` AS `shared_external`,`r`.`locale` AS `locale`,`r`.`name_de` AS `name_de`,`r`.`name_en` AS `name_en`,`r`.`name_fr` AS `name_fr`,`r`.`description_de` AS `description_de`,`r`.`description_en` AS `description_en`,`r`.`description_fr` AS `description_fr`,`r`.`servings_count` AS `servings_count`,`r`.`source_description_de` AS `source_description_de`,`r`.`source_description_en` AS `source_description_en`,`r`.`source_description_fr` AS `source_description_fr`,`r`.`source_url` AS `source_url`,`r`.`created` AS `created`,`r`.`modified` AS `modified`,`r`.`published` AS `published`,`r`.`difficulty` AS `difficulty`,`r`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,ifnull(`p`.`picture_id`,0) AS `picture_id`,ifnull(`p`.`picture_sortindex`,0) AS `picture_sortindex`,ifnull(`p`.`picture_name`,'') AS `picture_name`,ifnull(`p`.`picture_description`,'') AS `picture_description`,ifnull(`p`.`picture_hash`,'') AS `picture_hash`,ifnull(`p`.`picture_filename`,'') AS `picture_filename`,ifnull(`p`.`picture_full_path`,'') AS `picture_full_path`,cast(ifnull(`p`.`picture_uploaded`,'2000-01-01 00:00:00') as datetime) AS `picture_uploaded`,ifnull(`p`.`picture_width`,0) AS `picture_width`,ifnull(`p`.`picture_height`,0) AS `picture_height`,`rv`.`views` AS `views`,`rc`.`cooked` AS `cooked`,`v`.`votesum` AS `votesum`,`v`.`votes` AS `votes`,`v`.`avgvotes` AS `avgvotes`,`rr`.`votesum` AS `ratesum`,`rr`.`votes` AS `ratings`,`rr`.`avgvotes` AS `avgratings`,`s`.`stepscount` AS `stepscount`,ifnull(`s`.`preparationtime`,-1) AS `preparationtime`,ifnull(`s`.`cookingtime`,-1) AS `cookingtime`,ifnull(`s`.`chilltime`,-1) AS `chilltime` from ((((((`recipes` `r` left join `recipe_pictures` `p` on(`p`.`recipe_id` = `r`.`recipe_id` and `p`.`picture_sortindex` = 0)) join `voting_cooked` `rc` on(`rc`.`recipe_id` = `r`.`recipe_id`)) join `voting_views` `rv` on(`rv`.`recipe_id` = `r`.`recipe_id`)) join `voting_hearts` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) join `voting_difficulty` `rr` on(`rr`.`recipe_id` = `r`.`recipe_id`)) join `allrecipessteps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`,`r`.`user_id`,`r`.`edit_user_id`,`r`.`aigenerated`,`r`.`localized`,`r`.`placeholder`,`r`.`shared_internal`,`r`.`shared_external`,`r`.`locale`,`r`.`name_de`,`r`.`name_en`,`r`.`name_fr`,`r`.`description_de`,`r`.`description_en`,`r`.`description_fr`,`r`.`servings_count`,`r`.`source_description_de`,`r`.`source_description_en`,`r`.`source_description_fr`,`r`.`source_url`,`r`.`created`,`r`.`modified`,`r`.`published`,`r`.`difficulty`,`r`.`ingredientsGroupByStep`,`p`.`picture_id`,`p`.`picture_sortindex`,`p`.`picture_name`,`p`.`picture_description`,`p`.`picture_hash`,`p`.`picture_filename`,`p`.`picture_full_path`,`p`.`picture_uploaded`,`p`.`picture_width`,`p`.`picture_height`,`rv`.`views`,`rc`.`cooked`,`v`.`votesum`,`v`.`votes`,`v`.`avgvotes`,`s`.`stepscount`,`s`.`preparationtime`,`s`.`cookingtime`,`s`.`chilltime`;

ALTER ALGORITHM = UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `allrecipes_nouser` AS select `allrecipes`.`recipe_id` AS `recipe_id`,`allrecipes`.`user_id` AS `user_id`,`allrecipes`.`edit_user_id` AS `edit_user_id`,`allrecipes`.`aigenerated` AS `aigenerated`,`allrecipes`.`localized` AS `localized`,`allrecipes`.`placeholder` AS `placeholder`,`allrecipes`.`shared_internal` AS `shared_internal`,`allrecipes`.`shared_external` AS `shared_external`,`allrecipes`.`locale` AS `locale`,`allrecipes`.`name_de` AS `name_de`,`allrecipes`.`name_en` AS `name_en`,`allrecipes`.`name_fr` AS `name_fr`,`allrecipes`.`description_de` AS `description_de`,`allrecipes`.`description_en` AS `description_en`,`allrecipes`.`description_fr` AS `description_fr`,`allrecipes`.`servings_count` AS `servings_count`,`allrecipes`.`source_description_de` AS `source_description_de`,`allrecipes`.`source_description_en` AS `source_description_en`,`allrecipes`.`source_description_fr` AS `source_description_fr`,`allrecipes`.`source_url` AS `source_url`,`allrecipes`.`created` AS `created`,`allrecipes`.`modified` AS `modified`,`allrecipes`.`published` AS `published`,`allrecipes`.`difficulty` AS `difficulty`,`allrecipes`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,`allrecipes`.`picture_id` AS `picture_id`,`allrecipes`.`picture_sortindex` AS `picture_sortindex`,`allrecipes`.`picture_name` AS `picture_name`,`allrecipes`.`picture_description` AS `picture_description`,`allrecipes`.`picture_hash` AS `picture_hash`,`allrecipes`.`picture_filename` AS `picture_filename`,`allrecipes`.`picture_full_path` AS `picture_full_path`,`allrecipes`.`picture_uploaded` AS `picture_uploaded`,`allrecipes`.`picture_width` AS `picture_width`,`allrecipes`.`picture_height` AS `picture_height`,`allrecipes`.`views` AS `views`,`allrecipes`.`cooked` AS `cooked`,`allrecipes`.`votesum` AS `votesum`,`allrecipes`.`votes` AS `votes`,`allrecipes`.`avgvotes` AS `avgvotes`,`allrecipes`.`ratesum` AS `ratesum`,`allrecipes`.`ratings` AS `ratings`,`allrecipes`.`avgratings` AS `avgratings`,`allrecipes`.`stepscount` AS `stepscount`,`allrecipes`.`preparationtime` AS `preparationtime`,`allrecipes`.`cookingtime` AS `cookingtime`,`allrecipes`.`chilltime` AS `chilltime` from `allrecipes` where `allrecipes`.`shared_external` = 1;
```

## Update Categories

```sql
ALTER TABLE `categories`
    CHANGE COLUMN `techname` `name_de` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `catid`,
    ADD COLUMN `name_en` VARCHAR(64) NOT NULL AFTER `name_de`,
    ADD COLUMN `name_fr` VARCHAR(64) NOT NULL AFTER `name_en`;

UPDATE `categories` SET `name_de` = 'Regional', `name_en` = 'Regional', `name_fr` = 'Régional' WHERE `catid` = 1;
UPDATE `categories` SET `name_de` = 'Saisonal', `name_en` = 'Seasonal', `name_fr` = 'Saisonnier' WHERE `catid` = 2;
UPDATE `categories` SET `name_de` = 'Mahlzeit', `name_en` = 'Meal', `name_fr` = 'Repas' WHERE `catid` = 3;
UPDATE `categories` SET `name_de` = 'Zubereitung', `name_en` = 'Preparation', `name_fr` = 'Préparation' WHERE `catid` = 4;
UPDATE `categories` SET `name_de` = 'Anlass', `name_en` = 'Occasion', `name_fr` = 'Occasion' WHERE `catid` = 5;
UPDATE `categories` SET `name_de` = 'Hauptzutat', `name_en` = 'Key Ingredient', `name_fr` = 'Ingrédient clé' WHERE `catid` = 6;
UPDATE `categories` SET `name_de` = 'Ernährung', `name_en` = 'Diet', `name_fr` = 'Régime' WHERE `catid` = 7;
UPDATE `categories` SET `name_de` = 'Geschmack', `name_en` = 'Flavor', `name_fr` = 'Saveur' WHERE `catid` = 8;
UPDATE `categories` SET `name_de` = 'Dessert', `name_en` = 'Dessert', `name_fr` = 'Dessert' WHERE `catid` = 9;

ALTER TABLE `categoryitems`
    CHANGE COLUMN `techname` `name_de` VARCHAR(64) NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `catid`,
    ADD COLUMN `name_en` VARCHAR(64) NOT NULL AFTER `name_de`,
    ADD COLUMN `name_fr` VARCHAR(64) NOT NULL AFTER `name_en`;

UPDATE `categoryitems` SET `name_de` = 'Italienisch', `name_en` = 'Italian', `name_fr` = 'Italien' WHERE `itemid` = 15;
UPDATE `categoryitems` SET `name_de` = 'Spanisch', `name_en` = 'Spanish', `name_fr` = 'Espagnol' WHERE `itemid` = 16;
UPDATE `categoryitems` SET `name_de` = 'Französisch', `name_en` = 'French', `name_fr` = 'Français' WHERE `itemid` = 17;
UPDATE `categoryitems` SET `name_de` = 'Schwäbisch', `name_en` = 'Swabian', `name_fr` = 'Souabe' WHERE `itemid` = 18;
UPDATE `categoryitems` SET `name_de` = 'Asiatisch', `name_en` = 'Asian', `name_fr` = 'Asiatique' WHERE `itemid` = 37;
UPDATE `categoryitems` SET `name_de` = 'Mexikanisch', `name_en` = 'Mexican', `name_fr` = 'Mexicain' WHERE `itemid` = 38;
UPDATE `categoryitems` SET `name_de` = 'Indisch', `name_en` = 'Indian', `name_fr` = 'Indien' WHERE `itemid` = 39;
UPDATE `categoryitems` SET `name_de` = 'Orientalisch', `name_en` = 'Oriental', `name_fr` = 'Oriental' WHERE `itemid` = 40;
UPDATE `categoryitems` SET `name_de` = 'Frühling', `name_en` = 'Spring', `name_fr` = 'Printemps' WHERE `itemid` = 46;
UPDATE `categoryitems` SET `name_de` = 'Sommer', `name_en` = 'Summer', `name_fr` = 'Été' WHERE `itemid` = 47;
UPDATE `categoryitems` SET `name_de` = 'Herbst', `name_en` = 'Autumn', `name_fr` = 'Automne' WHERE `itemid` = 48;
UPDATE `categoryitems` SET `name_de` = 'Winter', `name_en` = 'Winter', `name_fr` = 'Hiver' WHERE `itemid` = 49;
UPDATE `categoryitems` SET `name_de` = 'Spargelsaison', `name_en` = 'Asparagus Season', `name_fr` = 'Saison des asperges' WHERE `itemid` = 50;
UPDATE `categoryitems` SET `name_de` = 'Brunch', `name_en` = 'Brunch', `name_fr` = 'Brunch' WHERE `itemid` = 1;
UPDATE `categoryitems` SET `name_de` = 'Frühstück', `name_en` = 'Breakfast', `name_fr` = 'Petit-déjeuner' WHERE `itemid` = 2;
UPDATE `categoryitems` SET `name_de` = 'Mittagessen', `name_en` = 'Lunch', `name_fr` = 'Déjeuner' WHERE `itemid` = 3;
UPDATE `categoryitems` SET `name_de` = 'Abendessen', `name_en` = 'Dinner', `name_fr` = 'Dîner' WHERE `itemid` = 4;
UPDATE `categoryitems` SET `name_de` = 'Kaffee und Kuchen', `name_en` = 'Coffee and Cake', `name_fr` = 'Café et gâteau' WHERE `itemid` = 5;
UPDATE `categoryitems` SET `name_de` = 'Teatime', `name_en` = 'Tea Time', `name_fr` = 'L\'heure du thé' WHERE `itemid` = 6;
UPDATE `categoryitems` SET `name_de` = 'Snack', `name_en` = 'Snack', `name_fr` = 'Snack' WHERE `itemid` = 7;
UPDATE `categoryitems` SET `name_de` = 'Beilage', `name_en` = 'Side Dish', `name_fr` = 'Accompagnement' WHERE `itemid` = 14;
UPDATE `categoryitems` SET `name_de` = 'Grillen', `name_en` = 'Grill', `name_fr` = 'Griller' WHERE `itemid` = 26;
UPDATE `categoryitems` SET `name_de` = 'Braten', `name_en` = 'Roast', `name_fr` = 'Rôtir' WHERE `itemid` = 27;
UPDATE `categoryitems` SET `name_de` = 'Backen', `name_en` = 'Bake', `name_fr` = 'Cuire au four' WHERE `itemid` = 28;
UPDATE `categoryitems` SET `name_de` = 'Dämpfen', `name_en` = 'Steam', `name_fr` = 'Cuire à la vapeur' WHERE `itemid` = 29;
UPDATE `categoryitems` SET `name_de` = 'Kochen', `name_en` = 'Boil', `name_fr` = 'Bouillir' WHERE `itemid` = 30;
UPDATE `categoryitems` SET `name_de` = 'Frittieren', `name_en` = 'Deep Fry', `name_fr` = 'Frire' WHERE `itemid` = 31;
UPDATE `categoryitems` SET `name_de` = 'Party', `name_en` = 'Party', `name_fr` = 'Fête' WHERE `itemid` = 8;
UPDATE `categoryitems` SET `name_de` = 'Weihnachten', `name_en` = 'Christmas', `name_fr` = 'Noël' WHERE `itemid` = 32;
UPDATE `categoryitems` SET `name_de` = 'Geburtstag', `name_en` = 'Birthday', `name_fr` = 'Anniversaire' WHERE `itemid` = 33;
UPDATE `categoryitems` SET `name_de` = 'Grillen', `name_en` = 'Barbecue', `name_fr` = 'Barbecue' WHERE `itemid` = 34;
UPDATE `categoryitems` SET `name_de` = 'Picknick', `name_en` = 'Picnic', `name_fr` = 'Pique-nique' WHERE `itemid` = 35;
UPDATE `categoryitems` SET `name_de` = 'Halloween', `name_en` = 'Halloween', `name_fr` = 'Halloween' WHERE `itemid` = 36;
UPDATE `categoryitems` SET `name_de` = 'Schwein', `name_en` = 'Pork', `name_fr` = 'Porc' WHERE `itemid` = 19;
UPDATE `categoryitems` SET `name_de` = 'Rind', `name_en` = 'Beef', `name_fr` = 'Bœuf' WHERE `itemid` = 20;
UPDATE `categoryitems` SET `name_de` = 'Huhn', `name_en` = 'Chicken', `name_fr` = 'Poulet' WHERE `itemid` = 21;
UPDATE `categoryitems` SET `name_de` = 'Gemüse', `name_en` = 'Vegetables', `name_fr` = 'Légumes' WHERE `itemid` = 22;
UPDATE `categoryitems` SET `name_de` = 'Pasta', `name_en` = 'Pasta', `name_fr` = 'Pâtes' WHERE `itemid` = 23;
UPDATE `categoryitems` SET `name_de` = 'Reis', `name_en` = 'Rice', `name_fr` = 'Riz' WHERE `itemid` = 24;
UPDATE `categoryitems` SET `name_de` = 'Vegetarisch', `name_en` = 'Vegetarian', `name_fr` = 'Végétarien' WHERE `itemid` = 10;
UPDATE `categoryitems` SET `name_de` = 'Vegan', `name_en` = 'Vegan', `name_fr` = 'Vegan' WHERE `itemid` = 11;
UPDATE `categoryitems` SET `name_de` = 'Glutenfrei', `name_en` = 'Gluten-Free', `name_fr` = 'Sans gluten' WHERE `itemid` = 25;
UPDATE `categoryitems` SET `name_de` = 'Scharf', `name_en` = 'Spicy', `name_fr` = 'Épicé' WHERE `itemid` = 12;
UPDATE `categoryitems` SET `name_de` = 'Fruchtig', `name_en` = 'Fruity', `name_fr` = 'Fruité' WHERE `itemid` = 13;
UPDATE `categoryitems` SET `name_de` = 'Kuchen', `name_en` = 'Cake', `name_fr` = 'Gâteau' WHERE `itemid` = 41;
UPDATE `categoryitems` SET `name_de` = 'Torte', `name_en` = 'Pie', `name_fr` = 'Tarte' WHERE `itemid` = 42;
UPDATE `categoryitems` SET `name_de` = 'Keks', `name_en` = 'Biscuit', `name_fr` = 'Biscuit' WHERE `itemid` = 43;
UPDATE `categoryitems` SET `name_de` = 'Eiscreme', `name_en` = 'Ice Cream', `name_fr` = 'Glace' WHERE `itemid` = 44;
UPDATE `categoryitems` SET `name_de` = 'Pudding', `name_en` = 'Pudding', `name_fr` = 'Pudding' WHERE `itemid` = 45;

ALTER ALGORITHM = UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `categoryitemsview` AS select `c`.`catid` AS `cat_id`,`c`.`name_de` AS `cat_name_de`,`c`.`name_en` AS `cat_name_en`,`c`.`name_fr` AS `cat_name_fr`,`c`.`icon` AS `cat_icon`,`c`.`modified` AS `cat_modified`,`i`.`itemid` AS `item_id`,`i`.`name_de` AS `item_name_de`,`i`.`name_en` AS `item_name_en`,`i`.`name_fr` AS `item_name_fr`,`i`.`icon` AS `item_icon`,`i`.`modified` AS `item_modified` from (`categoryitems` `i` join `categories` `c` on(`c`.`catid` = `i`.`catid`));
```

## Logging Changes

```sql
DELETE FROM `apilog`;

ALTER TABLE `apilog`
    AUTO_INCREMENT=1,
    ADD COLUMN `reporter` SET('Server','Client') NOT NULL AFTER `severity`,
    CHANGE COLUMN `host` `host` VARCHAR(128) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `reporter`,
    ADD COLUMN `agent` VARCHAR(256) NOT NULL DEFAULT '' AFTER `host`,
    CHANGE COLUMN `request_type` `request_type` VARCHAR(12) NOT NULL DEFAULT '' COMMENT 'GET, POST, etc' COLLATE 'utf8mb4_general_ci' AFTER `agent`,
    CHANGE COLUMN `request_uri` `request_uri` VARCHAR(256) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `request_type`,
    CHANGE COLUMN `request_length` `request_length` BIGINT(20) NOT NULL DEFAULT 0 AFTER `request_uri`,
    CHANGE COLUMN `message` `message` VARCHAR(1024) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `request_length`,
    DROP COLUMN `phpclass`,
    DROP COLUMN `phpmethod`,
    DROP COLUMN `phpline`,
    DROP COLUMN `payload`,
    ADD INDEX `host` (`host`),
    ADD INDEX `agent` (`agent`),
    ADD INDEX `request_uri` (`request_uri`),
    ADD INDEX `severity` (`severity`);

ALTER TABLE `apilog`
    CHANGE COLUMN `severity` `severity` SET('I','W','E') NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `when`;

```
