CREATE DATABASE netmap;
CREATE EXTENSION IF NOT EXISTS timescaledb;

\c netmap

CREATE TABLE graphs (
  time        TIMESTAMPTZ       NOT NULL,
  reporter    TEXT              NOT NULL,
  vertices    JSON              NULL,
  edges       JSON              NULL,
  packets     INTEGER           NULL
);

SELECT create_hypertable('graphs', 'time');