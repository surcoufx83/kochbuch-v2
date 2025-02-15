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
