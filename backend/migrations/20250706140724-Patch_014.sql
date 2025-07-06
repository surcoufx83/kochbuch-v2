
-- +migrate Up
ALTER TABLE `apilog`
	CHANGE COLUMN `message` `message` TEXT NOT NULL COLLATE 'utf8mb4_general_ci' AFTER `request_length`,
	ADD COLUMN `recipe_id` INT NULL DEFAULT NULL AFTER `message`;

-- +migrate Down
ALTER TABLE `apilog`
	CHANGE COLUMN `message` `message` VARCHAR(1024) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci' AFTER `request_length`,
	DROP COLUMN `recipe_id`;
