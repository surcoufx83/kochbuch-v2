
-- +migrate Up
DELETE FROM `recipe_categories`;

ALTER TABLE `recipe_categories`
	DROP PRIMARY KEY,
	ADD PRIMARY KEY (`recipe_id`, `catitem_id`) USING BTREE;
-- +migrate Down
