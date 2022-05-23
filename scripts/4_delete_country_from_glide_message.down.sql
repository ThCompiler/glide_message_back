ALTER TABLE glide_message
    ALTER COLUMN country SET DEFAULT 'Russian',
    ALTER COLUMN country SET not null,
    DROP COLUMN age;