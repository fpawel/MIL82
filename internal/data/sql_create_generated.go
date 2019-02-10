package data

const SQLCreate = `
PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';
--E:\Program Data\Аналитприбор\mil82\mil82.sqlite
--C:\Users\fpawel\AppData\Roaming\Аналитприбор\mil82\mil82.sqlite

CREATE TABLE IF NOT EXISTS party
(
  party_id          INTEGER PRIMARY KEY NOT NULL,
  created_at        TIMESTAMP           NOT NULL DEFAULT (datetime('now')),
  product_type      TEXT                NOT NULL DEFAULT '00.01',
  pgs1              REAL                NOT NULL DEFAULT 0 CHECK (pgs1 >= 0),
  pgs2              REAL                NOT NULL DEFAULT 50 CHECK (pgs2 >= 0),
  pgs3              REAL                NOT NULL DEFAULT 100 CHECK (pgs3 >= 0),
  pgs4              REAL                NOT NULL DEFAULT 100 CHECK (pgs4 >= 0),
  temperature_norm  REAL                NOT NULL DEFAULT 20,
  temperature_minus REAL                NOT NULL DEFAULT -60,
  temperature_plus  REAL                NOT NULL DEFAULT 80
);

CREATE TABLE IF NOT EXISTS product
(
  product_id INTEGER PRIMARY KEY NOT NULL,
  party_id   INTEGER             NOT NULL,
  serial     INTEGER             NOT NULL CHECK (serial > 0 ),
  place      INTEGER             NOT NULL CHECK (place >= 0),
  addr       INTEGER             NOT NULL CHECK (addr > 0),
  production BOOLEAN             NOT NULL CHECK (production IN (0, 1)) DEFAULT 0,
  CONSTRAINT unique_party_place UNIQUE (party_id, place),
  CONSTRAINT unique_party_serial UNIQUE (party_id, serial),
  CONSTRAINT unique_party_addr UNIQUE (party_id, addr),
  FOREIGN KEY (party_id) REFERENCES party (party_id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_test_error
(
  product_test_error_id INTEGER PRIMARY KEY NOT NULL,
  product_id       INTEGER NOT NULL,
  test_error_point TEXT    NOT NULL CHECK ( test_error_point IN
                                            ('norm', 'minus', 'plus', '90', 'norm2', 'tex1', 'tex2') ),
  gas_point        INTEGER NOT NULL CHECK ( gas_point >= 0 AND gas_point <= 4 ),
  value            REAL    NOT NULL,
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE,
  UNIQUE (product_id, test_error_point, gas_point)
);

CREATE TABLE IF NOT EXISTS product_coefficient
(
  product_coefficient_id INTEGER PRIMARY KEY NOT NULL,
  product_id  INTEGER NOT NULL,
  coefficient INTEGER NOT NULL CHECK ( coefficient >= 0 ),
  value       REAL    NOT NULL,
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE,
  UNIQUE (product_id, coefficient)
);

CREATE TABLE IF NOT EXISTS product_lin
(
  product_lin_id INTEGER PRIMARY KEY NOT NULL,
  product_id INTEGER NOT NULL,
  lin_point  INTEGER NOT NULL CHECK ( lin_point >= 0 AND lin_point <= 4),
  value      REAL    NOT NULL,
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE,
  UNIQUE (product_id, lin_point)
);

CREATE TABLE IF NOT EXISTS product_temp
(
  product_coefficient_id INTEGER PRIMARY KEY NOT NULL,
  product_id        INTEGER NOT NULL,
  gas_point         INTEGER NOT NULL CHECK ( gas_point >= 0 AND gas_point <= 4 ),
  temperature_point INTEGER NOT NULL CHECK ( temperature_point >= 0 AND temperature_point <= 3 ),
  temperature       REAL    NOT NULL,
  value             REAL    NOT NULL,
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE,
  UNIQUE (product_id, gas_point, temperature_point)
);

CREATE TABLE IF NOT EXISTS work
(
  work_id    INTEGER   NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL UNIQUE DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
  name       TEXT      NOT NULL
);

CREATE TABLE IF NOT EXISTS entry
(
  entry_id   INTEGER   NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL UNIQUE DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
  message    TEXT      NOT NULL,
  level      TEXT      NOT NULL,
  work_id    INTEGER   NOT NULL,
  FOREIGN KEY (work_id) REFERENCES work (work_id) ON DELETE CASCADE
);

CREATE VIEW IF NOT EXISTS work_info AS
SELECT *,
       EXISTS(
           SELECT level IN ('panic', 'error', 'fatal') FROM entry WHERE entry.work_id = work.work_id) AS error_occurred,
       CAST(STRFTIME('%Y', DATETIME(created_at, '+3 hours')) AS INTEGER)                              AS year,
       CAST(STRFTIME('%m', DATETIME(created_at, '+3 hours')) AS INTEGER)                              AS month,
       CAST(STRFTIME('%d', DATETIME(created_at, '+3 hours')) AS INTEGER)                              AS day
FROM work;`
