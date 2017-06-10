-- Config

create table if not exists config (
  key varchar not null unique check (key <> ''),
  value text not null
);

insert into config (key, value) values
  ('webl.title', 'Web Log'),
  ('webl.description', 'A place to speak your thoughts'),
  ('webl.baseurl', 'http://localhost:8080'),
  ('webl.jwt.secret', 'thirds-and-fifths');

create index if not exists config_key_idx on config(key);

--
--
