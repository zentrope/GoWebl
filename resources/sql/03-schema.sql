-- Update accounts to use a name and UUID

-- This will break if user is not super user, so, might have
-- to do this as a separate, non-automated step.
create extension if not exists "pgcrypto";

-- Add a unique 'uuid' column to author.

alter table author rename column handle to name;
alter table author rename constraint author_handle_key to author_name_key;
alter table author add column uuid uuid not null unique default gen_random_uuid();

-- Change the post.author column to uuid and set as fkey to author.uuid

alter table post drop constraint post_author_fkey;
update post set author=author.uuid from author where post.author=author.name;
alter table post alter column author set data type uuid using author::uuid;
alter table post rename column author to author_uuid;
alter table post add constraint post_author_uuid_fkey foreign key (author_uuid) references author(uuid);
