--migrator:up
CREATE TABLE public.users (
  "id"              uuid                        PRIMARY KEY DEFAULT gen_random_uuid(),
  "login"           text                        UNIQUE NOT NULL,
  "password_digest" text                        NOT NULL,
  "locale"          text                        NOT NULL DEFAULT 'en',
  "roles"           text[]                      NOT NULL DEFAULT '{"user"}',

  "created_at"      timestamp without time zone NOT NULL DEFAULT current_timestamp,
  "updated_at"      timestamp without time zone NOT NULL DEFAULT current_timestamp
);

--migrator:down
DROP TABLE public.users;
