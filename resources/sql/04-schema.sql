-- Add a date_published to posts.
--

alter table post add column
  date_published timestamp with time zone default current_timestamp;

update post set date_published=date_created;
