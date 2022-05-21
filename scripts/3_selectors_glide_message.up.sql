CREATE TABLE glide_message_countries
(
    id            bigserial not null primary key,
    county        citext    not null not null references countries (country_name) on delete cascade,
    glide_message bigserial not null references glide_message (id) on delete cascade
);

CREATE TABLE glide_message_languages
(
    id            bigserial not null primary key,
    language citext    not null references languages (language) on delete cascade,
    glide_message bigserial not null references glide_message (id) on delete cascade
);
