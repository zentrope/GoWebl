-- database
--

create type author_type as enum ('admin', 'user');
create type author_status as enum ('active', 'inactive');

create table if not exists author (
  handle varchar not null unique,
  email varchar not null unique,
  password varchar not null,
  status author_status not null default 'active',
  type author_type not null default 'user'
);

insert into author (email, handle, password)
  values ('keith@zentrope.com', 'keith', 'test1234');
insert into author (email, handle, password)
  values ('xan@zentrope.com', 'xan', 'test1234');

create type post_status as enum ('published', 'draft');

create table if not exists post (
  id serial primary key,
  author varchar references author(handle) not null,
  date_created timestamp with time zone default current_timestamp,
  date_updated timestamp with time zone default current_timestamp,
  status post_status not null default 'draft',
  text text default ''
);

insert into post (author, text)
  values ('keith', 'This is the first post.');
insert into post (author, text)
  values ('xan', 'This is the second post.');

-- To view all the enums
create or replace view vw_enums as
select
  t.typname, e.enumlabel, e.enumsortorder
from
  pg_enum e
join
  pg_type t ON e.enumtypid = t.oid;
