
-- +migrate Up
UPDATE `kochbuch`.`units` SET `fractional`=1 WHERE  `unit_id` IN (5, 7, 8, 9, 53, 55, 56, 60);

-- +migrate Down
UPDATE `kochbuch`.`units` SET `fractional`=0 WHERE  `unit_id` IN (5, 7, 8, 9, 53, 55, 56, 60);
