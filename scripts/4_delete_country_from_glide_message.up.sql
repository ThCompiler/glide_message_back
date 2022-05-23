ALTER TABLE glide_message
    ALTER COLUMN country DROP not null,
    ADD COLUMN age int default 0 not null;
