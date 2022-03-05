-- Request collection

create table if not exists request (
  address inet not null default '0.0.0.0',
  date_recorded timestamp with time zone default current_timestamp,
  method varchar not null default 'GET',
  path varchar not null,
  user_agent varchar not null,
  referer varchar not null default ''
);

create index if not exists request_address_idx on request (address);
create index if not exists request_path_idx on request (lower(path));
create index if not exists request_agent_idx on request (lower(user_agent));
