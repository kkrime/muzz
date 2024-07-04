DROP DATABASE IF EXISTS muzz;
CREATE DATABASE muzz;

\connect muzz
CREATE EXTENSION postgis;

CREATE TYPE public.gender AS ENUM (
	'M',
	'F'
);

-- TODO: add indexes

CREATE TABLE public.users(
  id SERIAL NOT NULL PRIMARY KEY,
	first_name CHARACTER VARYING(100) NOT NULL,
	last_name CHARACTER VARYING(100) NOT NULL,
	email CHARACTER VARYING(255) NOT NULL,
	password CHARACTER VARYING(60) NOT NULL,
  gender public.gender NOT NULL,
	dob DATE NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE public.login(
  id SERIAL NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES public.users(id),
  location GEOGRAPHY(POINT, 4326) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE public.swipe(
  id SERIAL NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES public.users(id),
  their_user_id INTEGER NOT NULL REFERENCES public.users(id),
  swipe_right BOOLEAN NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);

-- CREATE INDEX name_city ON public.ports(name);
-- CREATE INDEX name_city ON public.ports(primary_unloc);
-- CREATE UNIQUE INDEX no_duplicate_code ON public.ports(primary_unloc,deleted_at)
--    WHERE deleted_at IS null;

-- CREATE TABLE public.alias(
--     port_id INTEGER REFERENCES public.ports NOT NULL,
--     name CHARACTER VARYING(100) NOT NULL
-- );
-- CREATE UNIQUE INDEX no_duplicate_alias ON public.alias(port_id,name);

-- CREATE TABLE public.regions(
--     port_id INTEGER REFERENCES public.ports NOT NULL,
--     name CHARACTER VARYING(100) NOT NULL
-- );
-- CREATE UNIQUE INDEX no_duplicate_regions ON public.regions(port_id,name);

-- CREATE TABLE public.unlocs(
--     port_id INTEGER REFERENCES public.ports NOT NULL,
--     name CHARACTER VARYING(100) NOT NULL
-- );
-- CREATE UNIQUE INDEX no_duplicate_unlocs ON public.unlocs(port_id,name);
