CREATE DATABASE zomato_delivery;

\c zomato_delivery;

CREATE TABLE IF NOT EXISTS deliveries (
    delivery_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'PREPARING'
); 