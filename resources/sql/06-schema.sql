-- indexes
--

create index if not exists author_email_idx on author(lower(email));
create index if not exists post_date_published_idx on post(date_published);
create index if not exists post_author_uuid_idx on post(author_uuid);
