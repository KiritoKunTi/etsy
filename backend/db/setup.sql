drop table if exists product_photo;
drop table if exists product_likes;
drop table if exists product_comments;
drop table if exists sessions;
drop table if exists messages;
drop table if exists chats;
drop table if exists product_parameters;
drop table if exists products;
drop table if exists categories;
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
    is_active     boolean,
    created_at    timestamp not null
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id) on delete cascade,
  created_at timestamp not null   
);
create table categories(
    id              serial primary key,
    name            varchar(244),
    amo_products    integer
);
create table products(
    id              serial primary key,
    user_id         integer references users(id) on delete cascade,
    category_id     integer references categories(id) on delete cascade,
    name            varchar(255),
    photo           varchar(350),
    price           integer,
    amount          integer,
    description     varchar,
    amo_likes       integer,
    amo_comments    integer,
    amo_ratings     integer,
    rating          decimal,
    is_active       boolean,
    created_at      timestamp
);
create table product_parameters(
    id          serial primary key,
    product_id  integer references products(id) on delete cascade,
    key         varchar(255),
    value       varchar
);

create table product_photo(
    id          serial primary key,
    product_id  integer references products(id) on delete cascade,
    photo       varchar(320)
);

create table product_likes(
    id         serial primary key,
    product_id integer references products(id) on delete cascade,
    user_id    integer references users(id) on delete cascade,
    created_at timestamp
);

create table product_comments(
    id          serial primary key,
    product_id  integer references products(id) on delete cascade,
    user_id     integer references users(id) on delete cascade,
    text        varchar,
    created_at  timestamp
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
create or replace function product_insertion_category_function()
    returns trigger
    as $$
    begin
    update categories set amo_products = amo_products + 1 where id = NEW.category_id;
    return new;
    end;
    $$
    language 'plpgsql';

create or replace trigger product_insertion_category
    after insert on
    products for each row
    execute procedure product_insertion_category_function();


create or replace function product_delete_category_function()
    returns trigger
    as $$
begin
update categories set amo_products = amo_products - 1 where id = OLD.category_id;
return new;
end;
    $$
language 'plpgsql';

create or replace trigger product_delete_category
    after delete on
    products for each row
    execute procedure product_insertion_category_function();
alter table users alter "photo" set default 'private/avatar/default.jpg';
alter table users alter "language_code" set default 'en';
alter table users alter "description" set default '';
alter table users alter "is_active" set default true;
alter table products alter "amo_likes" set default 0;
alter table products alter "amo_comments" set default 0;
alter table products alter "amo_ratings" set default 0;
alter table categories alter "amo_products" set default 0;
alter table products alter "rating" set default 0;
alter table products alter "photo" set default 'private/product/default.jpg';
alter table products alter "is_active" set default true;
insert into categories(name) values
                                  ('Electronics'),
                                  (E'Men\'s Fashion');


