CREATE TABLE orders (
  order_uid       TEXT PRIMARY KEY,
  track_number    TEXT NOT NULL,
  entry           TEXT,
  locale          TEXT,
  internal_signature TEXT,
  customer_id     TEXT,
  delivery_service TEXT,
  shardkey        TEXT,
  sm_id           INTEGER,
  date_created    TIMESTAMPTZ,
  oof_shard       TEXT,
  created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE delivery (
  order_uid  TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
  name       TEXT,
  phone      TEXT,
  zip        TEXT,
  city       TEXT,
  address    TEXT,
  region     TEXT,
  email      TEXT,
  PRIMARY KEY (order_uid)
);

CREATE TABLE payment (
  order_uid     TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
  transaction   TEXT,
  request_id    TEXT,
  currency      TEXT,
  provider      TEXT,
  amount        BIGINT,
  payment_dt    BIGINT,          -- unix seconds (как в примере)
  bank          TEXT,
  delivery_cost BIGINT,
  goods_total   BIGINT,
  custom_fee    BIGINT,
  PRIMARY KEY (order_uid)
);

CREATE TABLE items (
  id            BIGSERIAL PRIMARY KEY,
  order_uid     TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id       BIGINT,
  track_number  TEXT,
  price         BIGINT,
  rid           TEXT,
  name          TEXT,
  sale          INTEGER,
  size          TEXT,
  total_price   BIGINT,
  nm_id         BIGINT,
  brand         TEXT,
  status        INTEGER
);

CREATE INDEX idx_items_order_uid ON items(order_uid);
