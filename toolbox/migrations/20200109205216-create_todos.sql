--migrator:up
CREATE TABLE public.todos (
  "id"              uuid                        PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id"         uuid                        REFERENCES public.users (id) NOT NULL,
  "body"            text                        NOT NULL,

  "created_at"      timestamp without time zone NOT NULL DEFAULT current_timestamp,
  "updated_at"      timestamp without time zone NOT NULL DEFAULT current_timestamp
);

--migrator:down
DROP TABLE public.todos;
