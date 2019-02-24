DO $$
  BEGIN
    CREATE TYPE SCALE_GAS_TYPE AS ENUM ( 'gas1', 'gas2', 'gas3');
    CREATE TYPE LIN_GAS_TYPE AS ENUM ( 'gas1', 'gas2', 'gas3', 'gas4', 'gas5');
    CREATE TYPE SCALE_TEMPERATURE_TYPE AS ENUM ( 'temperature_minus', 'temperature_norm', 'temperature_plus', 'temperature90');
    CREATE TYPE PRODUCT_TEST_TYPE AS ENUM (
      'test_norm', 'test_minus', 'test_plus', 'test90', 'test_norm2', 'test2', 'test3');
  EXCEPTION
    WHEN duplicate_object THEN null;
  END $$;

CREATE TABLE IF NOT EXISTS party
(
  party_id          SERIAL PRIMARY KEY       NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  product_type      TEXT                     NOT NULL DEFAULT '00.01',
  pgs1              REAL                     NOT NULL DEFAULT 0 CHECK (pgs1 >= 0),
  pgs2              REAL                     NOT NULL DEFAULT 50 CHECK (pgs2 >= 0),
  pgs3              REAL                     NOT NULL DEFAULT 100 CHECK (pgs3 >= 0),
  pgs_lin_12        REAL                     NOT NULL DEFAULT 25 CHECK (pgs_lin_12 >= 0),
  pgs_lin_22        REAL                     NOT NULL DEFAULT 25 CHECK (pgs_lin_22 >= 0),
  temperature_norm  REAL                     NOT NULL DEFAULT 20,
  temperature_minus REAL                     NOT NULL DEFAULT -60,
  temperature_plus  REAL                     NOT NULL DEFAULT 80
);

CREATE TABLE IF NOT EXISTS product
(
  product_id    SERIAL PRIMARY KEY NOT NULL,
  party_id      INTEGER            NOT NULL,
  serial_number SMALLINT           NOT NULL CHECK (serial_number > 0 ),
  place         SMALLINT           NOT NULL CHECK (place >= 0),
  addr          SMALLINT           NOT NULL CHECK (addr > 0),
  production    BOOLEAN            NOT NULL DEFAULT FALSE,
  UNIQUE (party_id, addr),
  UNIQUE (party_id, place),
  UNIQUE (party_id, serial_number),
  FOREIGN KEY (party_id) REFERENCES party (party_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_checkup
(
  product_id INTEGER NOT NULL,
  test       PRODUCT_TEST_TYPE,
  gas        SCALE_GAS_TYPE,
  value      REAL    NOT NULL,
  PRIMARY KEY (product_id, test, gas),
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_coefficient
(
  product_id  INTEGER NOT NULL,
  coefficient INTEGER NOT NULL CHECK ( coefficient >= 0 ),
  value       REAL    NOT NULL,
  PRIMARY KEY (product_id, coefficient),
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_lin
(
  product_id INTEGER      NOT NULL,
  gas        LIN_GAS_TYPE NOT NULL,
  value      REAL         NOT NULL,
  PRIMARY KEY (product_id, gas),
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_temperature_comp
(
  product_id        INTEGER                NOT NULL,
  gas               SCALE_GAS_TYPE         NOT NULL,
  temperature       SCALE_TEMPERATURE_TYPE NOT NULL,
  temperature_value REAL                   NOT NULL,
  value             REAL                   NOT NULL,
  PRIMARY KEY (product_id, temperature, gas),
  FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);
