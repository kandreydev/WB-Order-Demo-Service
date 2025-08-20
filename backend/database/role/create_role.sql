CREATE ROLE orders_user WITH LOGIN PASSWORD 'orders_pass';
CREATE DATABASE orders_db OWNER orders_user;
