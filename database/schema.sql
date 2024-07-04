DROP DATABASE IF EXISTS muzz;
CREATE DATABASE muzz;

\connect muzz
CREATE EXTENSION postgis;

CREATE TYPE public.gender AS ENUM (
	'M',
	'F'
);

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
CREATE INDEX user_email_idx ON public.users(email)
   WHERE deleted_at IS null;
CREATE INDEX user_email_all_idx ON public.users(email);
CREATE INDEX user_gender_idx ON public.users(gender);

CREATE TABLE public.logins(
  id SERIAL NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES public.users(id),
  location GEOGRAPHY(POINT, 4326) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX user_id_login_idx ON public.logins(user_id);
CREATE INDEX created_at_login_idx ON public.logins(created_at);

CREATE TABLE public.swipes(
  id SERIAL NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES public.users(id),
  their_user_id INTEGER NOT NULL REFERENCES public.users(id),
  swipe_right BOOLEAN NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX user_swipe_idx ON public.swipes(user_id, their_user_id,swipe_righT);
CREATE INDEX attractive_idx ON public.swipes(their_user_id, swipe_right);

-- CREATE INDEX name_city ON public.ports(name);
-- CREATE INDEX name_city ON public.ports(primary_unloc);
-- CREATE UNIQUE INDEX no_duplicate_code ON public.ports(primary_unloc,deleted_at)
--    WHERE deleted_at IS null;
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
