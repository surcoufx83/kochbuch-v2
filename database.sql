-- Exportiere Struktur von Tabelle kochbuch.activities_tracking
CREATE TABLE IF NOT EXISTS `activities_tracking` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `type` enum('show','create','update','delete') NOT NULL,
  `when` timestamp NOT NULL DEFAULT current_timestamp(),
  `user_id` mediumint(8) unsigned DEFAULT NULL,
  `recipe_id` int(10) unsigned DEFAULT NULL,
  `recipe_ai_id` mediumint(8) unsigned DEFAULT NULL,
  `step_id` bigint(20) unsigned DEFAULT NULL,
  `ingredient_id` bigint(20) unsigned DEFAULT NULL,
  `unit_id` smallint(5) unsigned DEFAULT NULL,
  `picture_id` int(10) unsigned DEFAULT NULL,
  `category_id` tinyint(3) unsigned DEFAULT NULL,
  `other_type` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=856 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.apicache
CREATE TABLE IF NOT EXISTS `apicache` (
  `route` varchar(512) NOT NULL,
  `ts` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `response` text NOT NULL,
  PRIMARY KEY (`route`,`ts`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.apilog
CREATE TABLE IF NOT EXISTS `apilog` (
  `entry` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `when` datetime NOT NULL DEFAULT current_timestamp(),
  `severity` set('I','W','E','F') NOT NULL,
  `host` varchar(128) NOT NULL,
  `request_length` bigint(20) NOT NULL,
  `request_uri` varchar(256) NOT NULL,
  `request_type` varchar(256) NOT NULL,
  `phpclass` varchar(128) NOT NULL,
  `phpmethod` varchar(128) NOT NULL,
  `phpline` int(11) NOT NULL,
  `message` text NOT NULL,
  `payload` text DEFAULT NULL,
  PRIMARY KEY (`entry`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2618 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.categories
CREATE TABLE IF NOT EXISTS `categories` (
  `catid` tinyint(3) unsigned NOT NULL AUTO_INCREMENT,
  `techname` varchar(32) NOT NULL,
  `icon` varchar(64) NOT NULL DEFAULT '',
  `modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`catid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.categoryitems
CREATE TABLE IF NOT EXISTS `categoryitems` (
  `itemid` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `catid` tinyint(3) unsigned NOT NULL,
  `techname` varchar(32) NOT NULL,
  `icon` varchar(64) NOT NULL DEFAULT '',
  `modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`itemid`) USING BTREE,
  KEY `FK__categories` (`catid`) USING BTREE,
  CONSTRAINT `FK__categories` FOREIGN KEY (`catid`) REFERENCES `categories` (`catid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipes
CREATE TABLE IF NOT EXISTS `recipes` (
  `recipe_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` mediumint(8) unsigned DEFAULT NULL,
  `edit_user_id` mediumint(8) unsigned DEFAULT NULL,
  `aigenerated` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `localized` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `recipe_placeholder` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `recipe_public_internal` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `recipe_public_external` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `recipe_name` varchar(256) NOT NULL,
  `recipe_name_de` varchar(256) DEFAULT NULL,
  `recipe_name_en` varchar(256) DEFAULT NULL,
  `recipe_description` varchar(1024) NOT NULL,
  `recipe_description_de` varchar(1024) DEFAULT NULL,
  `recipe_description_en` varchar(1024) DEFAULT NULL,
  `recipe_eater` tinyint(3) unsigned NOT NULL,
  `recipe_source_desc` varchar(1024) DEFAULT '',
  `recipe_source_url` varchar(256) DEFAULT '',
  `recipe_created` timestamp NOT NULL DEFAULT current_timestamp(),
  `recipe_edited` timestamp NULL DEFAULT NULL,
  `recipe_modified` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `recipe_published` timestamp NULL DEFAULT NULL,
  `recipe_difficulty` enum('0','1','2','3') NOT NULL DEFAULT '0',
  `ingredientsGroupByStep` tinyint(1) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`recipe_id`),
  KEY `FK_recipes_users` (`user_id`),
  CONSTRAINT `FK_recipes_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=75 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_ai_scanner
CREATE TABLE IF NOT EXISTS `recipe_ai_scanner` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `userid` mediumint(8) unsigned NOT NULL,
  `request` text NOT NULL,
  `response_raw` text NOT NULL,
  `recipeid` int(10) unsigned DEFAULT NULL,
  `tokens_in` int(10) unsigned DEFAULT NULL,
  `token_out` int(10) unsigned DEFAULT NULL,
  `when` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_categories
CREATE TABLE IF NOT EXISTS `recipe_categories` (
  `recipe_id` int(10) unsigned NOT NULL,
  `catitem_id` smallint(5) unsigned NOT NULL,
  `user_id` mediumint(8) unsigned NOT NULL,
  PRIMARY KEY (`recipe_id`,`catitem_id`,`user_id`) USING BTREE,
  KEY `FK_recipe_categories_categoryitems` (`catitem_id`) USING BTREE,
  KEY `recipe_id` (`recipe_id`) USING BTREE,
  KEY `FK_recipe_categories_users` (`user_id`) USING BTREE,
  CONSTRAINT `FK_recipe_categories_categoryitems` FOREIGN KEY (`catitem_id`) REFERENCES `categoryitems` (`itemid`) ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_categories_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_categories_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_ingredients
CREATE TABLE IF NOT EXISTS `recipe_ingredients` (
  `ingredient_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `recipe_id` int(10) unsigned NOT NULL,
  `step_id` bigint(20) unsigned DEFAULT NULL,
  `sortindex` smallint(5) unsigned NOT NULL,
  `unit_id` smallint(5) unsigned DEFAULT NULL,
  `ingredient_quantity` decimal(10,3) DEFAULT NULL,
  `ingredient_description` varchar(128) NOT NULL,
  `ingredient_description_de` varchar(128) DEFAULT NULL,
  `ingredient_description_en` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`ingredient_id`) USING BTREE,
  KEY `FK_recipe_ingredients_recipes` (`recipe_id`),
  KEY `FK_recipe_ingredients_units` (`unit_id`),
  KEY `FK_recipe_ingredients_recipe_steps` (`step_id`),
  CONSTRAINT `FK_recipe_ingredients_recipe_steps` FOREIGN KEY (`step_id`) REFERENCES `recipe_steps` (`step_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_ingredients_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_ingredients_units` FOREIGN KEY (`unit_id`) REFERENCES `units` (`unit_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=912 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_pictures
CREATE TABLE IF NOT EXISTS `recipe_pictures` (
  `picture_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `recipe_id` int(10) unsigned NOT NULL,
  `user_id` mediumint(8) unsigned DEFAULT NULL,
  `picture_sortindex` tinyint(3) unsigned NOT NULL,
  `picture_name` varchar(128) NOT NULL,
  `picture_description` varchar(1024) NOT NULL,
  `picture_hash` varchar(32) NOT NULL,
  `picture_filename` varchar(256) NOT NULL,
  `picture_full_path` varchar(1024) NOT NULL,
  `picture_uploaded` timestamp NOT NULL DEFAULT current_timestamp(),
  `picture_width` int(11) NOT NULL,
  `picture_height` int(11) NOT NULL,
  PRIMARY KEY (`picture_id`),
  UNIQUE KEY `recipe_id_picture_sortindex` (`recipe_id`,`picture_sortindex`),
  KEY `FK_recipe_pictures_users` (`user_id`),
  CONSTRAINT `FK_recipe_pictures_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_pictures_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=120 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_steps
CREATE TABLE IF NOT EXISTS `recipe_steps` (
  `step_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `recipe_id` int(10) unsigned NOT NULL,
  `step_no` tinyint(3) unsigned NOT NULL,
  `step_title` varchar(128) NOT NULL,
  `step_title_de` varchar(128) DEFAULT NULL,
  `step_title_en` varchar(128) DEFAULT NULL,
  `step_data` text NOT NULL,
  `step_data_de` text DEFAULT NULL,
  `step_data_en` text DEFAULT NULL,
  `step_time_preparation` int(11) DEFAULT NULL,
  `step_time_cooking` int(11) DEFAULT NULL,
  `step_time_chill` int(11) DEFAULT NULL,
  PRIMARY KEY (`step_id`),
  UNIQUE KEY `recipe_id_step_no` (`recipe_id`,`step_no`),
  CONSTRAINT `FK_recipe_steps_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=268 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_voting_cooked
CREATE TABLE IF NOT EXISTS `recipe_voting_cooked` (
  `recipe_id` int(10) unsigned NOT NULL,
  `user_id` mediumint(8) unsigned DEFAULT NULL,
  `when` datetime NOT NULL DEFAULT current_timestamp(),
  KEY `FK_recipe_cooked_recipes` (`recipe_id`) USING BTREE,
  KEY `FK_recipe_cooked_users` (`user_id`) USING BTREE,
  CONSTRAINT `FK_recipe_cooked_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_cooked_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_voting_difficulty
CREATE TABLE IF NOT EXISTS `recipe_voting_difficulty` (
  `recipe_id` int(10) unsigned NOT NULL,
  `user_id` mediumint(8) unsigned DEFAULT NULL,
  `when` datetime NOT NULL DEFAULT current_timestamp(),
  `value` tinyint(3) unsigned NOT NULL,
  KEY `recipe_id` (`recipe_id`) USING BTREE,
  KEY `FK_recipe_ratings_users` (`user_id`) USING BTREE,
  CONSTRAINT `FK_recipe_ratings_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_ratings_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.recipe_voting_hearts
CREATE TABLE IF NOT EXISTS `recipe_voting_hearts` (
  `recipe_id` int(10) unsigned NOT NULL,
  `user_id` mediumint(8) unsigned DEFAULT NULL,
  `when` datetime NOT NULL DEFAULT current_timestamp(),
  `value` tinyint(3) unsigned NOT NULL,
  KEY `FK_recipe_vote_users` (`user_id`) USING BTREE,
  KEY `recipe_id` (`recipe_id`) USING BTREE,
  CONSTRAINT `FK_recipe_vote_recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_recipe_vote_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.units
CREATE TABLE IF NOT EXISTS `units` (
  `unit_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `useunit_id` smallint(5) unsigned DEFAULT NULL,
  `localized` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `unit_name` varchar(16) NOT NULL,
  `unit_name_de` varchar(64) DEFAULT NULL,
  `unit_name_en` varchar(64) DEFAULT NULL,
  `created` datetime NOT NULL DEFAULT current_timestamp(),
  `updated` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`unit_id`),
  UNIQUE KEY `unit_name` (`unit_name`),
  KEY `FK_units_units` (`useunit_id`),
  CONSTRAINT `FK_units_units` FOREIGN KEY (`useunit_id`) REFERENCES `units` (`unit_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=92 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.users
CREATE TABLE IF NOT EXISTS `users` (
  `user_id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) DEFAULT NULL,
  `oauth_user_name` varchar(32) DEFAULT NULL,
  `user_firstname` varchar(64) NOT NULL DEFAULT '',
  `user_lastname` varchar(64) NOT NULL DEFAULT '',
  `user_fullname` varchar(128) NOT NULL DEFAULT '',
  `user_hash` varchar(32) DEFAULT NULL,
  `user_isactivated` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Wurde das Konto durch einen Admin freigeschaltet',
  `user_isadmin` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `user_password` varchar(256) NOT NULL DEFAULT '',
  `user_email` varchar(256) NOT NULL,
  `user_email_validation` varchar(256) DEFAULT NULL,
  `user_email_validated` datetime DEFAULT NULL,
  `user_last_activity` datetime DEFAULT NULL,
  `user_avatar` varchar(40) DEFAULT NULL,
  `user_registration_completed` datetime DEFAULT NULL,
  `user_adconsent` datetime DEFAULT NULL,
  `user_betatester` tinyint(3) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_email` (`user_email`),
  UNIQUE KEY `user_name` (`user_name`),
  UNIQUE KEY `oauth_user_name` (`oauth_user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.user_logins
CREATE TABLE IF NOT EXISTS `user_logins` (
  `login_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` mediumint(8) unsigned NOT NULL,
  `login_time` timestamp NOT NULL DEFAULT current_timestamp(),
  `login_token` varchar(128) NOT NULL,
  `login_password` varchar(128) NOT NULL,
  `login_keep` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `login_oauthdata` text DEFAULT NULL,
  PRIMARY KEY (`login_id`),
  KEY `user_id` (`user_id`),
  KEY `login_time` (`login_time`),
  KEY `login_keep` (`login_keep`),
  CONSTRAINT `FK_user_logins_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Tabelle kochbuch.user_tokens
CREATE TABLE IF NOT EXISTS `user_tokens` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `userid` mediumint(8) unsigned DEFAULT NULL,
  `page` varchar(64) NOT NULL,
  `access_token` varchar(256) NOT NULL,
  `refresh_token` varchar(256) NOT NULL,
  `hash` varchar(64) NOT NULL,
  `fingerprint` varchar(64) DEFAULT NULL,
  `expires_in` int(11) NOT NULL,
  `created` datetime NOT NULL DEFAULT current_timestamp(),
  `valid` datetime NOT NULL,
  `lastaccess` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `FK_user_tokens_users` (`userid`) USING BTREE,
  CONSTRAINT `FK_user_tokens_users` FOREIGN KEY (`userid`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Exportiere Struktur von Trigger kochbuch.categories_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `categories_after_delete` AFTER DELETE ON `categories` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `category_id`) VALUES('delete', OLD.catid);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.categories_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `categories_after_insert` AFTER INSERT ON `categories` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `category_id`) VALUES('create', NEW.catid);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipes_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipes_after_delete` AFTER DELETE ON `recipes` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`) VALUES('delete', OLD.user_id, OLD.recipe_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipes_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipes_after_insert` AFTER INSERT ON `recipes` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`) VALUES('create', NEW.user_id, NEW.recipe_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipes_after_update
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipes_after_update` AFTER UPDATE ON `recipes` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`) VALUES('update', NEW.edit_user_id, NEW.recipe_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_categories_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_categories_after_delete` AFTER DELETE ON `recipe_categories` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `category_id`) VALUES('delete', OLD.recipe_id, OLD.catitem_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_categories_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_categories_after_insert` AFTER INSERT ON `recipe_categories` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `category_id`) VALUES('create', NEW.recipe_id, NEW.catitem_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_ingredients_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_ingredients_after_delete` AFTER DELETE ON `recipe_ingredients` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `step_id`, `ingredient_id`) VALUES('delete', OLD.recipe_id, OLD.step_id, OLD.ingredient_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_ingredients_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_ingredients_after_insert` AFTER INSERT ON `recipe_ingredients` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `step_id`, `ingredient_id`) VALUES('create', NEW.recipe_id, NEW.step_id, NEW.ingredient_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_pictures_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_pictures_after_delete` AFTER DELETE ON `recipe_pictures` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`, `picture_id`) VALUES('delete', OLD.user_id, OLD.recipe_id, OLD.picture_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_pictures_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_pictures_after_insert` AFTER INSERT ON `recipe_pictures` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`, `picture_id`) VALUES('create', NEW.user_id, NEW.recipe_id, NEW.picture_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_steps_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_steps_after_delete` AFTER DELETE ON `recipe_steps` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `step_id`) VALUES('delete', OLD.recipe_id, OLD.step_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_steps_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_steps_after_insert` AFTER INSERT ON `recipe_steps` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `step_id`) VALUES('create', NEW.recipe_id, NEW.step_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_voting_cooked_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_voting_cooked_after_delete` AFTER DELETE ON `recipe_voting_cooked` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `other_type`) VALUES('delete', OLD.recipe_id, 'cooked');
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_voting_cooked_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_voting_cooked_after_insert` AFTER INSERT ON `recipe_voting_cooked` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `other_type`) VALUES('create', NEW.recipe_id, 'cooked');
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_voting_difficulty_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_voting_difficulty_after_delete` AFTER DELETE ON `recipe_voting_difficulty` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `other_type`) VALUES('delete', OLD.recipe_id, 'difficulty');
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_voting_difficulty_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_voting_difficulty_after_insert` AFTER INSERT ON `recipe_voting_difficulty` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `other_type`) VALUES('create', NEW.recipe_id, 'difficulty');
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_voting_hearts_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_voting_hearts_after_delete` AFTER DELETE ON `recipe_voting_hearts` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `other_type`) VALUES('delete', OLD.recipe_id, 'hearts');
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.recipe_voting_hearts_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `recipe_voting_hearts_after_insert` AFTER INSERT ON `recipe_voting_hearts` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `recipe_id`, `other_type`) VALUES('create', NEW.recipe_id, 'hearts');
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.units_after_delete
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `units_after_delete` AFTER DELETE ON `units` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `unit_id`) VALUES('delete', OLD.unit_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Exportiere Struktur von Trigger kochbuch.units_after_insert
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
DELIMITER //
CREATE TRIGGER `units_after_insert` AFTER INSERT ON `units` FOR EACH ROW BEGIN
INSERT INTO activities_tracking(`type`, `unit_id`) VALUES('create', NEW.unit_id);
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `allrecipes` AS select `r`.`recipe_id` AS `recipe_id`,`r`.`user_id` AS `user_id`,`r`.`edit_user_id` AS `edit_user_id`,`r`.`aigenerated` AS `aigenerated`,`r`.`localized` AS `localized`,`r`.`recipe_placeholder` AS `recipe_placeholder`,`r`.`recipe_public_internal` AS `recipe_public_internal`,`r`.`recipe_public_external` AS `recipe_public_external`,`r`.`recipe_name` AS `recipe_name`,`r`.`recipe_name_de` AS `recipe_name_de`,`r`.`recipe_name_en` AS `recipe_name_en`,`r`.`recipe_description` AS `recipe_description`,`r`.`recipe_description_de` AS `recipe_description_de`,`r`.`recipe_description_en` AS `recipe_description_en`,`r`.`recipe_eater` AS `recipe_eater`,`r`.`recipe_source_desc` AS `recipe_source_desc`,`r`.`recipe_source_url` AS `recipe_source_url`,`r`.`recipe_created` AS `recipe_created`,`r`.`recipe_edited` AS `recipe_edited`,`r`.`recipe_modified` AS `recipe_modified`,`r`.`recipe_published` AS `recipe_published`,`r`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,`p`.`picture_id` AS `picture_id`,`p`.`picture_sortindex` AS `picture_sortindex`,`p`.`picture_name` AS `picture_name`,`p`.`picture_description` AS `picture_description`,`p`.`picture_hash` AS `picture_hash`,`p`.`picture_filename` AS `picture_filename`,`p`.`picture_full_path` AS `picture_full_path`,`p`.`picture_uploaded` AS `picture_uploaded`,`p`.`picture_width` AS `picture_width`,`p`.`picture_height` AS `picture_height`,`rv`.`views` AS `views`,`rc`.`cooked` AS `cooked`,`v`.`votesum` AS `votesum`,`v`.`votes` AS `votes`,`v`.`avgvotes` AS `avgvotes`,`rr`.`votesum` AS `ratesum`,`rr`.`votes` AS `ratings`,`rr`.`avgvotes` AS `avgratings`,`s`.`stepscount` AS `stepscount`,`s`.`preparationtime` AS `preparationtime`,`s`.`cookingtime` AS `cookingtime`,`s`.`chilltime` AS `chilltime` from ((((((`recipes` `r` left join `recipe_pictures` `p` on(`p`.`recipe_id` = `r`.`recipe_id` and `p`.`picture_sortindex` = 0)) join `voting_cooked` `rc` on(`rc`.`recipe_id` = `r`.`recipe_id`)) join `voting_views` `rv` on(`rv`.`recipe_id` = `r`.`recipe_id`)) join `voting_hearts` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) join `voting_difficulty` `rr` on(`rr`.`recipe_id` = `r`.`recipe_id`)) join `allrecipessteps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`,`r`.`user_id`,`r`.`edit_user_id`,`r`.`aigenerated`,`r`.`localized`,`r`.`recipe_placeholder`,`r`.`recipe_public_internal`,`r`.`recipe_public_external`,`r`.`recipe_name`,`r`.`recipe_name_de`,`r`.`recipe_name_en`,`r`.`recipe_description`,`r`.`recipe_description_de`,`r`.`recipe_description_en`,`r`.`recipe_eater`,`r`.`recipe_source_desc`,`r`.`recipe_source_url`,`r`.`recipe_created`,`r`.`recipe_edited`,`r`.`recipe_modified`,`r`.`recipe_published`,`r`.`ingredientsGroupByStep`,`p`.`picture_id`,`p`.`picture_sortindex`,`p`.`picture_name`,`p`.`picture_description`,`p`.`picture_hash`,`p`.`picture_filename`,`p`.`picture_full_path`,`p`.`picture_uploaded`,`p`.`picture_width`,`p`.`picture_height`,`rv`.`views`,`rc`.`cooked`,`v`.`votesum`,`v`.`votes`,`v`.`avgvotes`,`s`.`stepscount`,`s`.`preparationtime`,`s`.`cookingtime`,`s`.`chilltime`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `allrecipessteps` AS select `r`.`recipe_id` AS `recipe_id`,count(`s`.`step_no`) AS `stepscount`,sum(if(ifnull(`s`.`step_time_preparation`,0) < 0,0,`s`.`step_time_preparation`)) AS `preparationtime`,sum(if(ifnull(`s`.`step_time_cooking`,0) < 0,0,`s`.`step_time_cooking`)) AS `cookingtime`,sum(if(ifnull(`s`.`step_time_chill`,0) < 0,0,`s`.`step_time_chill`)) AS `chilltime` from (`recipes` `r` left join `recipe_steps` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `allrecipes_nouser` AS select `allrecipes`.`recipe_id` AS `recipe_id`,`allrecipes`.`user_id` AS `user_id`,`allrecipes`.`edit_user_id` AS `edit_user_id`,`allrecipes`.`aigenerated` AS `aigenerated`,`allrecipes`.`localized` AS `localized`,`allrecipes`.`recipe_placeholder` AS `recipe_placeholder`,`allrecipes`.`recipe_public_internal` AS `recipe_public_internal`,`allrecipes`.`recipe_public_external` AS `recipe_public_external`,`allrecipes`.`recipe_name` AS `recipe_name`,`allrecipes`.`recipe_name_de` AS `recipe_name_de`,`allrecipes`.`recipe_name_en` AS `recipe_name_en`,`allrecipes`.`recipe_description` AS `recipe_description`,`allrecipes`.`recipe_description_de` AS `recipe_description_de`,`allrecipes`.`recipe_description_en` AS `recipe_description_en`,`allrecipes`.`recipe_eater` AS `recipe_eater`,`allrecipes`.`recipe_source_desc` AS `recipe_source_desc`,`allrecipes`.`recipe_source_url` AS `recipe_source_url`,`allrecipes`.`recipe_created` AS `recipe_created`,`allrecipes`.`recipe_edited` AS `recipe_edited`,`allrecipes`.`recipe_modified` AS `recipe_modified`,`allrecipes`.`recipe_published` AS `recipe_published`,`allrecipes`.`ingredientsGroupByStep` AS `ingredientsGroupByStep`,`allrecipes`.`picture_id` AS `picture_id`,`allrecipes`.`picture_sortindex` AS `picture_sortindex`,`allrecipes`.`picture_name` AS `picture_name`,`allrecipes`.`picture_description` AS `picture_description`,`allrecipes`.`picture_hash` AS `picture_hash`,`allrecipes`.`picture_filename` AS `picture_filename`,`allrecipes`.`picture_full_path` AS `picture_full_path`,`allrecipes`.`picture_uploaded` AS `picture_uploaded`,`allrecipes`.`picture_width` AS `picture_width`,`allrecipes`.`picture_height` AS `picture_height`,`allrecipes`.`views` AS `views`,`allrecipes`.`cooked` AS `cooked`,`allrecipes`.`votesum` AS `votesum`,`allrecipes`.`votes` AS `votes`,`allrecipes`.`avgvotes` AS `avgvotes`,`allrecipes`.`ratesum` AS `ratesum`,`allrecipes`.`ratings` AS `ratings`,`allrecipes`.`avgratings` AS `avgratings`,`allrecipes`.`stepscount` AS `stepscount`,`allrecipes`.`preparationtime` AS `preparationtime`,`allrecipes`.`cookingtime` AS `cookingtime`,`allrecipes`.`chilltime` AS `chilltime` from `allrecipes` where `allrecipes`.`recipe_public_external` = 1;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `allrecipetextdata` AS select `r`.`recipe_id` AS `recipe_id`,`r`.`recipe_public_internal` AS `recipe_public_internal`,`r`.`recipe_public_external` AS `recipe_public_external`,`r`.`recipe_name` AS `recipe_name`,`r`.`recipe_name_de` AS `recipe_name_de`,`r`.`recipe_name_en` AS `recipe_name_en`,`r`.`recipe_description` AS `recipe_description`,`r`.`recipe_description_de` AS `recipe_description_de`,`r`.`recipe_description_en` AS `recipe_description_en`,`i`.`recipe_ingredients` AS `recipe_ingredients`,`i`.`recipe_ingredients_de` AS `recipe_ingredients_de`,`i`.`recipe_ingredients_en` AS `recipe_ingredients_en`,`s`.`recipe_steps` AS `recipe_steps`,`s`.`recipe_steps_de` AS `recipe_steps_de`,`s`.`recipe_steps_en` AS `recipe_steps_en`,concat(`r`.`recipe_name`,'\n',`r`.`recipe_description`,'\n',`i`.`recipe_ingredients`,'\n',`s`.`recipe_steps`) AS `recipe_fulltext`,concat(`r`.`recipe_name_de`,'\n',`r`.`recipe_description_de`,'\n',`i`.`recipe_ingredients_de`,'\n',`s`.`recipe_steps_de`) AS `recipe_fulltext_de`,concat(`r`.`recipe_name_en`,'\n',`r`.`recipe_description_en`,'\n',`i`.`recipe_ingredients_en`,'\n',`s`.`recipe_steps_en`) AS `recipe_fulltext_en`,concat(`r`.`recipe_name`,'\n',`r`.`recipe_description`,'\n',`s`.`recipe_steps`) AS `recipe_fulltext_noingredients`,concat(`r`.`recipe_name_de`,'\n',`r`.`recipe_description_de`,'\n',`s`.`recipe_steps_de`) AS `recipe_fulltext_noingredients_de`,concat(`r`.`recipe_name_en`,'\n',`r`.`recipe_description_en`,'\n',`s`.`recipe_steps_en`) AS `recipe_fulltext_noingredients_en` from ((`recipes` `r` join `recipe_ingredients_view` `i` on(`i`.`recipe_id` = `r`.`recipe_id`)) join `recipe_steps_view` `s` on(`s`.`recipe_id` = `r`.`recipe_id`)) where (`r`.`recipe_public_internal` = 1 or `r`.`recipe_public_external` = 1) and `r`.`recipe_placeholder` = 0;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `apicache_view` AS select `a`.`route` AS `route`,`a`.`ts` AS `ts`,`a`.`response` AS `response` from (`apicache` `a` join `lastactivityview` `l` on(`l`.`ts` = `a`.`ts`));

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `categoryitemsview` AS select `i`.`itemid` AS `itemid`,`i`.`techname` AS `itemname`,`i`.`icon` AS `itemicon`,`i`.`modified` AS `itemmodified`,`c`.`catid` AS `catid`,`c`.`techname` AS `catname`,`c`.`icon` AS `caticon`,`c`.`modified` AS `catmodified` from (`categoryitems` `i` join `categories` `c` on(`c`.`catid` = `i`.`catid`));

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `lastactivityview` AS select max(`activities_tracking`.`when`) AS `ts` from `activities_tracking`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `recipecategoriesview` AS select `rc`.`recipe_id` AS `recipe_id`,`ci`.`catid` AS `catid`,`ci`.`itemid` AS `itemid`,count(`rc`.`user_id`) AS `votes` from (`recipe_categories` `rc` join `categoryitems` `ci` on(`ci`.`itemid` = `rc`.`catitem_id`)) group by `rc`.`recipe_id`,`ci`.`catid`,`ci`.`itemid`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `recipes_my` AS select `recipes`.`recipe_id` AS `recipe_id`,`recipes`.`user_id` AS `user_id`,`recipes`.`recipe_name` AS `recipe_name`,`recipes`.`recipe_public_internal` AS `recipe_public_internal`,`recipes`.`recipe_public_external` AS `recipe_public_external` from `recipes`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `recipe_ingredients_view` AS select `recipe_ingredients`.`recipe_id` AS `recipe_id`,group_concat(distinct `recipe_ingredients`.`ingredient_description` separator '\n') AS `recipe_ingredients`,group_concat(distinct `recipe_ingredients`.`ingredient_description_de` separator '\n') AS `recipe_ingredients_de`,group_concat(distinct `recipe_ingredients`.`ingredient_description_en` separator '\n') AS `recipe_ingredients_en` from `recipe_ingredients` group by `recipe_ingredients`.`recipe_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `recipe_steps_view` AS select `recipe_steps`.`recipe_id` AS `recipe_id`,group_concat(distinct `recipe_steps`.`step_data` separator '\n') AS `recipe_steps`,group_concat(distinct `recipe_steps`.`step_data_de` separator '\n') AS `recipe_steps_de`,group_concat(distinct `recipe_steps`.`step_data_en` separator '\n') AS `recipe_steps_en` from `recipe_steps` group by `recipe_steps`.`recipe_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `units_usage` AS select `u`.`unit_id` AS `unit_id`,`u`.`unit_name` AS `unit_name`,count(`i`.`ingredient_id`) AS `usage_count` from (`units` `u` left join `recipe_ingredients` `i` on(`i`.`unit_id` = `u`.`unit_id`)) group by `u`.`unit_id`,`u`.`unit_name`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `user_statistics` AS select `u`.`user_id` AS `user_id`,(select count(`rvv`.`recipe_id`) from `recipe_voting_views` `rvv` where `rvv`.`user_id` = `u`.`user_id`) AS `recipes_viewed`,(select count(distinct `rvv`.`recipe_id`) from `recipe_voting_views` `rvv` where `rvv`.`user_id` = `u`.`user_id`) AS `distinct_recipes_viewed`,(select count(`rvc`.`recipe_id`) from `recipe_voting_cooked` `rvc` where `rvc`.`user_id` = `u`.`user_id`) AS `recipes_cooked`,(select count(distinct `rvc`.`recipe_id`) from `recipe_voting_cooked` `rvc` where `rvc`.`user_id` = `u`.`user_id`) AS `distinct_recipes_cooked`,(select count(`rvd`.`recipe_id`) from `recipe_voting_difficulty` `rvd` where `rvd`.`user_id` = `u`.`user_id`) AS `recipes_voted_difficulty`,(select count(distinct `rvd`.`recipe_id`) from `recipe_voting_difficulty` `rvd` where `rvd`.`user_id` = `u`.`user_id`) AS `distinct_recipes_voted_difficulty`,(select count(`rvh`.`recipe_id`) from `recipe_voting_hearts` `rvh` where `rvh`.`user_id` = `u`.`user_id`) AS `recipes_voted_hearts`,(select count(distinct `rvh`.`recipe_id`) from `recipe_voting_hearts` `rvh` where `rvh`.`user_id` = `u`.`user_id`) AS `distinct_recipes_voted_hearts`,(select count(distinct `r`.`recipe_id`) from `recipes` `r` where `r`.`user_id` = `u`.`user_id` and `r`.`recipe_placeholder` = 0) AS `recipes_created`,(select count(distinct `r`.`recipe_id`) from `recipes` `r` where `r`.`user_id` = `u`.`user_id` and `r`.`recipe_placeholder` = 0 and `r`.`aigenerated` = 1) AS `recipes_aigenerated`,(select count(distinct `r`.`recipe_id`) from `recipes` `r` where `r`.`user_id` = `u`.`user_id` and `r`.`recipe_placeholder` = 0 and (`r`.`recipe_public_internal` = 1 or `r`.`recipe_public_external` = 1)) AS `recipes_published`,(select count(distinct `r`.`recipe_id`) from `recipes` `r` where `r`.`user_id` = `u`.`user_id` and `r`.`recipe_placeholder` = 0 and `r`.`recipe_public_external` = 1) AS `recipes_published_external`,(select count(distinct `p`.`picture_id`) from `recipe_pictures` `p` where `p`.`user_id` = `u`.`user_id`) AS `recipes_pictures_uploaded` from `users` `u` group by `u`.`user_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `voting_cooked` AS select `r`.`recipe_id` AS `recipe_id`,count(`v`.`when`) AS `cooked` from (`recipes` `r` left join `recipe_voting_cooked` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `voting_difficulty` AS select `r`.`recipe_id` AS `recipe_id`,ifnull(sum(`v`.`value`),0) AS `votesum`,count(`v`.`when`) AS `votes`,ifnull(sum(`v`.`value`) / count(`v`.`when`),0) AS `avgvotes` from (`recipes` `r` left join `recipe_voting_difficulty` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `voting_hearts` AS select `r`.`recipe_id` AS `recipe_id`,ifnull(sum(`v`.`value`),0) AS `votesum`,count(`v`.`when`) AS `votes`,ifnull(sum(`v`.`value`) / count(`v`.`when`),0) AS `avgvotes` from (`recipes` `r` left join `recipe_voting_hearts` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`;

-- Entferne temporäre Tabelle und erstelle die eigentliche View
CREATE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `voting_views` AS select `r`.`recipe_id` AS `recipe_id`,count(`v`.`when`) AS `views` from (`recipes` `r` left join `recipe_voting_views` `v` on(`v`.`recipe_id` = `r`.`recipe_id`)) group by `r`.`recipe_id`;
