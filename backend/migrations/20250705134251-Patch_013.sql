-- +migrate Up
DROP VIEW IF EXISTS `allrecipetextdata`;

DROP VIEW IF EXISTS `lastactivityview`;

DROP VIEW IF EXISTS `recipecategoriesview`;

DROP VIEW IF EXISTS `recipe_ingredients_view`;

DROP VIEW IF EXISTS `recipe_steps_view`;

DROP VIEW IF EXISTS `units_usage`;

DROP VIEW IF EXISTS `user_statistics`;

-- +migrate Down