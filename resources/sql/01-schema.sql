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
  values ('root@example.com', 'root',
    '243261243130244676522f2f447452782e5036687451506b302f45684f65384a417a597853356e2f646a554237386358374e4c71635674574d47724f');

create type post_status as enum ('published', 'draft');

create table if not exists post (
  id serial primary key,
  uuid varchar not null unique,
  author varchar references author(handle) not null,
  date_created timestamp with time zone default current_timestamp,
  date_updated timestamp with time zone default current_timestamp,
  status post_status not null default 'draft',
  slugline text not null unique,
  text text default ''
);

insert into post (author, uuid, status, slugline, text) values ('root',
'b8ec300c-b0f1-4338-b2a7-ca5f06c1fe33',
'published',
'First Post',
'Welcome to my blog.');

-- To view all the enums
create or replace view vw_enums as
select
  t.typname, e.enumlabel, e.enumsortorder
from
  pg_enum e
join
  pg_type t ON e.enumtypid = t.oid;
