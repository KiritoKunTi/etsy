-- product
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
    execute procedure product_delete_category_function();


-- like
create or replace function like_insertion_product_function()
    returns trigger
    as $$
begin
update products set amo_likes = amo_likes + 1 where id = NEW.product_id;
return new;
end;
    $$
language 'plpgsql';

create or replace trigger like_insertion_product
    after insert on
    product_likes for each row
    execute procedure like_insertion_product_function();

create or replace function like_delete_product_function()
    returns trigger
    as $$
begin
update products set amo_likes = amo_likes - 1 where id = OLD.product_id;
return new;
end;
    $$
language 'plpgsql';

create or replace trigger like_delete_product
    after delete on
    product_likes for each row
    execute procedure like_delete_product_function();

-- comment

create or replace function comment_insertion_product_function()
    returns trigger
    as $$
begin
update products set amo_comments = amo_comments + 1 where id = NEW.product_id;
return new;
end;
    $$
language 'plpgsql';

create or replace trigger comment_delete_product
    after insert on
    product_comments for each row
    execute procedure comment_insertion_product_function();

create or replace function comment_delete_product_function()
    returns trigger
    as $$
begin
update products set amo_comments = amo_comments - 1 where id = OLD.product_id;
return new;
end;
    $$
language 'plpgsql';

create or replace trigger comment_insertion_product
    after delete on
    product_comments for each row
    execute procedure comment_delete_product_function();
