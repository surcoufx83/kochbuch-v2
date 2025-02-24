-- +migrate Up
ALTER TABLE `recipe_categories`
	ADD COLUMN `created` DATETIME NOT NULL DEFAULT current_timestamp() AFTER `user_id`;

-- +migrate Down
ALTER TABLE `recipe_categories`
	DROP COLUMN `created`;