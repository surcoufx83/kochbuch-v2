-- +migrate Up

-- +migrate StatementBegin
CREATE TRIGGER `recipe_pictures_after_update` AFTER UPDATE ON `recipe_pictures` FOR EACH ROW BEGIN
	UPDATE `recipes` SET `modified` = CURRENT_TIMESTAMP() WHERE `recipe_id` = NEW.recipe_id;
END
-- +migrate StatementEnd

DROP TRIGGER `recipe_pictures_after_insert`;

-- +migrate StatementBegin
CREATE TRIGGER `recipe_pictures_after_insert` AFTER INSERT ON `recipe_pictures` FOR EACH ROW BEGIN
    INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`, `picture_id`) VALUES('create', NEW.user_id, NEW.recipe_id, NEW.picture_id);
    UPDATE `recipes` SET `modified` = CURRENT_TIMESTAMP() WHERE `recipe_id` = NEW.recipe_id;
END
-- +migrate StatementEnd

DROP TRIGGER `recipe_pictures_after_delete`;

-- +migrate StatementBegin
CREATE TRIGGER `recipe_pictures_after_delete` AFTER DELETE ON `recipe_pictures` FOR EACH ROW BEGIN
    INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`, `picture_id`) VALUES('delete', OLD.user_id, OLD.recipe_id, OLD.picture_id);
    UPDATE `recipes` SET `modified` = CURRENT_TIMESTAMP() WHERE `recipe_id` = OLD.recipe_id;
END
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER `recipe_pictures_after_update`;
DROP TRIGGER `recipe_pictures_after_insert`;
DROP TRIGGER `recipe_pictures_after_delete`;

-- +migrate StatementBegin
CREATE TRIGGER `recipe_pictures_after_insert` AFTER INSERT ON `recipe_pictures` FOR EACH ROW BEGIN
    INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`, `picture_id`) VALUES('create', NEW.user_id, NEW.recipe_id, NEW.picture_id);
END
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TRIGGER `recipe_pictures_after_delete` AFTER DELETE ON `recipe_pictures` FOR EACH ROW BEGIN
    INSERT INTO activities_tracking(`type`, `user_id`, `recipe_id`, `picture_id`) VALUES('delete', OLD.user_id, OLD.recipe_id, OLD.picture_id);
END
-- +migrate StatementEnd