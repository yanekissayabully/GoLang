create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) unique not null, -- Добавили email, должен быть уникальным
    age int,                             -- Добавили возраст, может быть пустым
    created_at timestamp default now()   -- Добавили дату создания, ставится автоматически
);

-- Вставляем обновленные тестовые данные
insert into users (name, email, age) values
('John Doe', 'john.doe@example.com', 30),
('Jane Smith', 'jane.smith@example.com', 25);