CREATE DATABASE zomato_store;

\c zomato_store;

CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    items TEXT[] NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PREPARING'
); 