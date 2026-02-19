create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) unique not null,
    age int,
    created_at timestamp default now(),
    deleted_at timestamp null -- поле для мягкого удаления
);

insert into users (name, email, age) values
('John Doe', 'john.doe@example.com', 30),
('Jane Smith', 'jane.smith@example.com', 25),
('Bob Johnson', 'bob.johnson@example.com', 35);