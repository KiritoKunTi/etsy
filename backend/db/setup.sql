drop table if exists sessions;
drop table if exists messages;
drop table if exists chats;
drop table if exists users;
drop table if exists products;
drop table if exists product_parameters;
drop table if exists categories;

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
create table categories(
    id              serial primary key,
    name            varchar(244),
    amo_products    integer
);
create table products(
    id              serial primary key,
    user_id         integer references users(id),
    category_id     integer references categories(id),
    name            varchar(255),
    price           integer,
    amount          integer,
    description     varchar,
    amo_likes       integer,
    amo_comments    integer,
    amo_ratings     integer,
    rating          decimal,
    created_at      timestamp
);
create table product_parameters(
    id          serial primary key,
    product_id  integer references products(id),
    key         varchar(255),
    value       varchar
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
alter table users alter "photo" set default 'private/photo/default.jpg';
alter table users alter "language_code" set default 'en';
alter table users alter "description" set default '';
alter table products alter "amo_likes" set default 0;
alter table products alter "amo_comments" set default 0;
alter table product alter "amo_ratings" set default 0;
alter table categories alter "amo_products" set default 0;
insert into categories(name) values
                                  ("Electronics"),
                                  ("Men's Fashion");


