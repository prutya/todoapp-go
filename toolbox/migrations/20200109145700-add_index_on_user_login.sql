--migrator:up
CREATE INDEX index_on_users_login ON users (login);

--migrator:down
DROP INDEX index_on_users_login;
