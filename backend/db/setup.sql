drop table if exists sessions;
drop table if exists messages;
drop table if exists chats;
drop table if exists users;


create table users (
    id            serial primary key,
    uuid          varchar(64) not null unique,
    username      varchar(255) not null unique,
    email         varchar(255) not null unique,
    password      varchar(255) not null,
    first_name    varchar(255) not null,
    last_name     varchar(255),
    is_shop       boolean default false,
    photo         varchar(255),
    language_code varchar(255),
    description   text,
    created_at    timestamp not null
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id),
  created_at timestamp not null   
);

-- create table chats(
--     id              serial primary key,
--     user1_id        integer references users(id),
--     user2_id        integer references users(id),
--     last_message_id integer references  messages(id),
--     created_at      timestamp
-- );
-- create table messages(
--     id           serial primary key,
--     sender_id    integer references users(id),
--     chat_id      integer references chats(id),
--     type         varchar(255) not null,
--     document     varchar(255),
--     photo        varchar(255),
--     video        varchar(255),
--     voice        varchar(255),
--     text_message text,
--     created_at   timestamp
-- );
-- alter table users alter "photo" set default 'private/photo/default-avatar.jpg';
