SET client_encoding = 'UTF8';

CREATE TABLE account (
    uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    pass TEXT NOT NULL,
    email TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_account_uuid ON account(uuid);
INSERT INTO account (uuid, name, pass, email) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'Алеша', '1703', 'alesha@mail.ru');
INSERT INTO account (uuid, name, pass, email) 
    VALUES ('00000000-0000-0000-0000-000000000002', 'Степаша', '3223', 'killer@mail.ru');

CREATE TYPE cyrrency_code AS ENUM ('RUB', 'USD', 'EUR');

CREATE TABLE wallet (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id TEXT NOT NULL, 
    currency cyrrency_code NOT NULL, 
    count DOUBLE PRECISION,
    FOREIGN KEY (user_id) REFERENCES account(uuid) ON DELETE NO ACTION
);
ALTER TABLE wallet ADD CONSTRAINT uniq_wallet_user_cur UNIQUE (user_id, currency);
CREATE UNIQUE INDEX idx_wallet_id ON wallet(id);
INSERT INTO wallet (user_id, currency, count) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'EUR', 0);
INSERT INTO wallet (user_id, currency, count) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'USD', 17.3);
INSERT INTO wallet (user_id, currency, count) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'RUB', 10000.0);
    
CREATE TABLE cyrrency (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    note cyrrency_code NOT NULL,
    rate DOUBLE PRECISION
);
CREATE UNIQUE INDEX idx_cyrrency_id ON cyrrency(id);
INSERT INTO cyrrency (note, rate) VALUES ('RUB', 1.0);
INSERT INTO cyrrency (note, rate) VALUES ('USD', 105.0);
INSERT INTO cyrrency (note, rate) VALUES ('EUR', 110.0);
    