CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS languages
(
    language            citext       not null primary key,
    picture            text         not null,
);

CREATE TABLE IF NOT EXISTS countries
(
    country_name        citext       not null primary key,
    picture            text         not null,
);


CREATE TABLE IF NOT EXISTS users
(
    nickname            citext       not null primary key,
    fullname            text         not null,
    about               text      default ''   not null,
    password            text         not null,
    age                 int          not null,
    avatar              text     default ''    not null,
    country             citext   null foreign key countries(country_name) on delete cascade
);

CREATE TABLE IF NOT EXISTS user_language
(
    id   bigserial not null primary key,
    nickname        citext     no null foreign key users(nickname) on delete cascade,
    language        citext     no null foreign key languages(language) on delete cascade
);

CREATE TABLE IF NOT EXISTS chat
(
    id              bigserial not null primary key,
    author          citext      not null no null foreign key users(nickname) on delete cascade,
    companion          citext      not null no null foreign key users(nickname) on delete cascade
);

CREATE UNIQUE INDEX unique_companions_of_chat ON chat USING BTREE (hstore(ARRAY[author,companion], ARRAY[null,null]));

CREATE TABLE IF NOT EXISTS messages
(
    id   bigserial         not null primary key,
    message        text              not null,
    chat        bigint            not null references chat (id) on delete cascade,
    picture       text      default ''     not null,
    created        timestamptz default now()::timestamptz not null,
    author          citext      not null no null foreign key users(nickname) on delete cascade
);

CREATE TABLE IF NOT EXISTS glide_message
(
    id      bigserial         not null primary key,
    message        text              not null,
    title        text            not null ,
    picture       text      default ''     not null,
    created        timestamptz default now()::timestamptz not null,
    author          citext      not null no null foreign key users(nickname) on delete cascade,
    country          citext      not null no null foreign key countries(country) on delete cascade
);

CREATE TABLE IF NOT EXISTS glide_users
(
    id   bigserial not null primary key,
    visited_user        citext     no null foreign key users(nickname) on delete cascade,
    glide_message        bigserial     no null foreign key glide_message(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS black_list
(
    id              bigserial not null primary key,
    author          citext      not null no null foreign key users(nickname) on delete cascade,
    fobbiged_user   citext      not null no null foreign key users(nickname) on delete cascade
);
