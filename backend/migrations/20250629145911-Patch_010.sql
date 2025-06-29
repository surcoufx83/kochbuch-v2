
-- +migrate Up
DROP TRIGGER IF EXISTS `recipe_pictures_after_delete`;
DROP TRIGGER IF EXISTS `recipe_pictures_after_insert`;
DROP TRIGGER IF EXISTS `recipe_pictures_after_update`;
DROP TRIGGER IF EXISTS `recipe_ingredients_after_delete`;
DROP TRIGGER IF EXISTS `recipe_ingredients_after_insert`;
DROP TRIGGER IF EXISTS `recipe_steps_after_delete`;
DROP TRIGGER IF EXISTS `recipe_steps_after_insert`;
DROP TRIGGER IF EXISTS `categories_after_delete`;
DROP TRIGGER IF EXISTS `categories_after_insert`;
DROP TRIGGER IF EXISTS `recipes_after_delete`;
DROP TRIGGER IF EXISTS `recipes_after_insert`;
DROP TRIGGER IF EXISTS `recipes_after_update`;
DROP TRIGGER IF EXISTS `recipe_categories_after_delete`;
DROP TRIGGER IF EXISTS `recipe_categories_after_insert`;
DROP TRIGGER IF EXISTS `recipe_voting_cooked_after_delete`;
DROP TRIGGER IF EXISTS `recipe_voting_cooked_after_insert`;
DROP TRIGGER IF EXISTS `recipe_voting_difficulty_after_delete`;
DROP TRIGGER IF EXISTS `recipe_voting_difficulty_after_insert`;
DROP TRIGGER IF EXISTS `recipe_voting_hearts_after_delete`;
DROP TRIGGER IF EXISTS `recipe_voting_hearts_after_insert`;
DROP TRIGGER IF EXISTS `units_after_delete`;
DROP TRIGGER IF EXISTS `units_after_insert`;

-- +migrate Down
