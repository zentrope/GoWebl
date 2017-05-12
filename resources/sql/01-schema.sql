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
  values ('keith@example.com', 'keith',
    '243261243130244676522f2f447452782e5036687451506b302f45684f65384a417a597853356e2f646a554237386358374e4c71635674574d47724f');
insert into author (email, handle, password)
  values ('xan@example.com', 'xan',
    '24326124313024626a654965535441306c37756b6f7474496c6952757579773036556b7a5a76416445795735465855574175587a5451656155777147');

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

insert into post (author, uuid, status, slugline, text) values ('keith',
'b8ec300c-b0f1-4338-b2a7-ca5f06c1fe33',
'published',
'Xena of Amphipolous',
'# Xena of Amphipolous

Xena is a fictional character from Robert Tapert''s Xena: Warrior
Princess franchise. Co-created by Tapert and John Schulian, she first
appeared in the 1995–1999 television series Hercules: The Legendary
Journeys, before going on to appear in Xena: Warrior Princess TV show
and subsequent comic book of the same name. The Warrior Princess has
also appeared in the spin-off animated movie The Battle for Mount
Olympus, as well as numerous non-canon expanded universe material,
such as books and video games. Xena was played by New Zealand actress
Lucy Lawless.' );

insert into
       post (author, uuid, status, slugline, text) values ('xan',
'd08e9b3e-466b-423f-bd41-761f984bd0d0',
'published',
'Zeno of Elea',
'# Zeno of Elia

Little is known for certain about Zeno''s life. Although written
nearly a century after Zeno''s death, the primary source of
biographical information about Zeno is Plato''s Parmenides and he is
also mentioned in Aristotle''s Physics. In the dialogue of
Parmenides, Plato describes a visit to Athens by Zeno and Parmenides,
at a time when Parmenides is "about 65," Zeno is "nearly 40" and
Socrates is "a very young man". Assuming an age for Socrates of
around 20, and taking the date of Socrates'' birth as 469 BC gives an
approximate date of birth for Zeno of 490 BC."');

-- To view all the enums
create or replace view vw_enums as
select
  t.typname, e.enumlabel, e.enumsortorder
from
  pg_enum e
join
  pg_type t ON e.enumtypid = t.oid;
