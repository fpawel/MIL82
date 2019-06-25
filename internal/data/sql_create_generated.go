package data

const SQLCreate = `
PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS party
(
    party_id              INTEGER PRIMARY KEY NOT NULL,
    created_at            TIMESTAMP           NOT NULL DEFAULT (DATETIME('now', '+3 hours')),
    product_type          TEXT                NOT NULL DEFAULT '00.01',
    concentration_beg     REAL                NOT NULL DEFAULT 0 CHECK (concentration_beg >= 0),
    concentration_middle2 REAL                NOT NULL DEFAULT 25 CHECK (concentration_middle2 >= 0),
    concentration_middle  REAL                NOT NULL DEFAULT 50 CHECK (concentration_middle >= 0),
    concentration_end     REAL                NOT NULL DEFAULT 100 CHECK (concentration_end >= 0),
    temp_low              REAL                NOT NULL DEFAULT -60,
    temp_norm             REAL                NOT NULL DEFAULT 20,
    temp_high             REAL                NOT NULL DEFAULT 60,
    temp_90               REAL                NOT NULL DEFAULT 90
);

CREATE TABLE IF NOT EXISTS product
(
    product_id INTEGER PRIMARY KEY NOT NULL,
    party_id   INTEGER             NOT NULL,
    serial     SMALLINT            NOT NULL CHECK (serial > 0 ),
    addr       SMALLINT            NOT NULL CHECK (addr > 0),
    production BOOLEAN             NOT NULL DEFAULT FALSE,
    FOREIGN KEY (party_id) REFERENCES party (party_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_product_party_id_addr ON product (party_id, addr);
CREATE UNIQUE INDEX unique_product_party_id_serial ON product (party_id, serial);

CREATE TABLE IF NOT EXISTS product_temp_value
(
    product_id INTEGER NOT NULL,
    var     INTEGER NOT NULL CHECK (var >= 0),
    gas        TEXT    NOT NULL CHECK ( gas IN ('beg', 'mid', 'end') ),
    temp       TEXT    NOT NULL CHECK ( temp IN ('low', 'norm', 'high') ),
    value      REAL    NOT NULL,
    PRIMARY KEY (product_id, var, gas, temp),
    FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_work_value
(
    product_id INTEGER NOT NULL,
    var     INTEGER NOT NULL CHECK (var >= 0),
    work       TEXT    NOT NULL CHECK ( temp IN ('lin', 'checkup', 'tex1', 'tex2') ),
    gas        TEXT    NOT NULL CHECK ( gas IN ('gas1', 'gas2', 'gas3', 'gas4') ),
    temp       TEXT    NOT NULL CHECK ( temp IN ('low', 'norm', 'high', '90') ),
    value      REAL    NOT NULL,
    PRIMARY KEY (product_id, var, work, gas, temp),
    FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_coefficient
(
    product_id  INTEGER NOT NULL,
    coefficient INTEGER NOT NULL CHECK ( coefficient >= 0 ),
    value       REAL    NOT NULL,
    PRIMARY KEY (product_id, coefficient),
    FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);`
