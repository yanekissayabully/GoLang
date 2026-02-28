CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    genre VARCHAR(100) NOT NULL,
    rating FLOAT NOT NULL,
    description TEXT
);

-- Добавим начальные данные
INSERT INTO movies (title, genre, rating, description) VALUES
    ('The Shawshank Redemption', 'Drama', 9.3, 'Two imprisoned men bond over a number of years'),
    ('The Godfather', 'Crime', 9.2, 'The aging patriarch of an organized crime dynasty'),
    ('Pulp Fiction', 'Crime', 8.9, 'The lives of two mob hitmen, a boxer, and others intertwine');