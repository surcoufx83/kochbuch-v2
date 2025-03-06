
-- +migrate Up
ALTER TABLE `user_login_states`
	DROP INDEX `remoteaddr_useragent`;

-- +migrate Down
