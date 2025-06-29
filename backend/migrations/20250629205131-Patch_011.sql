
-- +migrate Up
CREATE TABLE `user_collections` (
	`collection_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` MEDIUMINT(8) UNSIGNED NOT NULL,
	`title` VARCHAR(256) NOT NULL COLLATE 'utf8mb4_general_ci',
	`description` TEXT NOT NULL COLLATE 'utf8mb4_general_ci',
	`created` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted` DATETIME NULL DEFAULT NULL,
	`published` DATETIME NULL DEFAULT NULL,
	PRIMARY KEY (`collection_id`) USING BTREE,
	INDEX `FK_user_collections_users` (`user_id`) USING BTREE,
	CONSTRAINT `FK_user_collections_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB;

INSERT INTO `user_collections`
	SELECT NULL, `user_id`, 'Mein Kochbuch', '', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), NULL, NULL
	FROM `users`;

CREATE TABLE `user_collection_recipes` (
	`collection_id` INT(10) UNSIGNED NOT NULL,
	`recipe_id` INT(10) UNSIGNED NOT NULL,
	`is_owner` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
	`remarks` TEXT NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
	`created` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	UNIQUE INDEX `collection_id_recipe_id` (`collection_id`, `recipe_id`) USING BTREE,
	INDEX `FK__recipes` (`recipe_id`) USING BTREE,
	CONSTRAINT `FK__recipes` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`recipe_id`) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT `FK__user_collections` FOREIGN KEY (`collection_id`) REFERENCES `user_collections` (`collection_id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB;

INSERT INTO `user_collection_recipes`(`collection_id`, `recipe_id`, `is_owner`)
	SELECT `collection_id`, `recipe_id`, 1
	FROM `user_collections`
	JOIN `recipes` ON `user_collections`.`user_id` = `recipes`.`user_id`;

-- +migrate Down
DROP TABLE IF EXISTS `user_collection_recipes`;
DROP TABLE IF EXISTS `user_collections`;