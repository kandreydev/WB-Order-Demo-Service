-- +migrate Down
DROP INDEX IF EXISTS idx_items_order_uid;

DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS orders;

DROP DATABASE IF EXISTS orders_db;
DROP ROLE IF EXISTS orders_user;
