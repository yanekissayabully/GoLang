create table if not exists users (
    id serial primary key,
    name varchar(255) not null
);

insert into users (name) values ('John Doe');