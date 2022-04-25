CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS hstore;

CREATE TABLE IF NOT EXISTS languages
(
    language citext not null primary key,
    picture  text   not null
);

CREATE TABLE IF NOT EXISTS countries
(
    country_name citext not null primary key,
    picture      text   not null
);


CREATE TABLE IF NOT EXISTS users
(
    nickname citext          not null primary key,
    fullname text            not null,
    about    text default '' not null,
    password text            not null,
    age      int             not null,
    avatar   text default '' not null,
    country  citext          null references countries (country_name) on delete cascade
);

CREATE TABLE IF NOT EXISTS user_language
(
    id       bigserial not null primary key,
    nickname citext    not null references users (nickname) on delete cascade,
    language citext    not null references languages (language) on delete cascade
);

CREATE TABLE IF NOT EXISTS chat
(
    id        bigserial not null primary key,
    author    citext    not null references users (nickname) on delete cascade,
    companion citext    not null references users (nickname) on delete cascade
);

CREATE UNIQUE INDEX unique_companions_of_chat ON chat USING BTREE (hstore(ARRAY [author, companion], ARRAY [null,null]));

CREATE TABLE IF NOT EXISTS messages
(
    id        bigserial                              not null primary key,
    message   text                                   not null,
    chat      bigint                                 not null references chat (id) on delete cascade,
    picture   text        default ''                 not null,
    created   timestamptz default now()::timestamptz not null,
    is_viewed bool        default false              not null,
    author    citext                                 not null not null references users (nickname) on delete cascade
);

CREATE TABLE IF NOT EXISTS glide_message
(
    id      bigserial                              not null primary key,
    message text                                   not null,
    title   text                                   not null,
    picture text        default ''                 not null,
    created timestamptz default now()::timestamptz not null,
    author  citext                                 not null not null references users (nickname) on delete cascade,
    country citext                                 not null not null references countries (country_name) on delete cascade
);

CREATE TABLE IF NOT EXISTS glide_users
(
    id            bigserial not null primary key,
    visited_user  citext    not null references users (nickname) on delete cascade,
    glide_message bigserial not null references glide_message (id) on delete cascade
);

CREATE TABLE IF NOT EXISTS black_list
(
    id            bigserial not null primary key,
    author        citext    not null not null references users (nickname) on delete cascade,
    fobbiged_user citext    not null not null references users (nickname) on delete cascade
);
