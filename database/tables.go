package database

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error{
	query := `
	
CREATE TABLE IF NOT EXISTS public.users
(
    id serial NOT NULL,
    username character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    role_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_username_key UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS public.roles
(
    id serial NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT roles_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.boards
(
    id serial NOT NULL,
    title character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    creator_id integer,
    is_public boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT boards_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.board_permissions
(
    id serial NOT NULL,
    board_id integer,
    user_id integer,
    can_edit boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT board_permissions_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.board_elements
(
    id serial NOT NULL,
    board_id integer,
    type character varying(50) COLLATE pg_catalog."default" NOT NULL,
    content text COLLATE pg_catalog."default",
    position_x integer,
    position_y integer,
    width integer,
    height integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT board_elements_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id)
    REFERENCES public.roles (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.boards
    ADD CONSTRAINT boards_creator_id_fkey FOREIGN KEY (creator_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.board_permissions
    ADD CONSTRAINT board_permissions_board_id_fkey FOREIGN KEY (board_id)
    REFERENCES public.boards (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.board_permissions
    ADD CONSTRAINT board_permissions_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.board_elements
    ADD CONSTRAINT board_elements_board_id_fkey FOREIGN KEY (board_id)
    REFERENCES public.boards (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

END;

	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}