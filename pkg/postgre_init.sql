SET client_encoding = 'UTF8';

CREATE TABLE account (
    uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    pass TEXT NOT NULL,
    email TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_account_uuid ON account(uuid);
INSERT INTO account (uuid, name, pass, email) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'Алеша', '$2a$10$LnNCNcub3VKVVGE.t4fBS.bL00/SydpIUeE1P2nKKEvrckhn8iNBq', 'alesha@mail.ru'); --pass = 1703
INSERT INTO account (uuid, name, pass, email) 
    VALUES ('00000000-0000-0000-0000-000000000002', 'Степаша', '$2a$10$zkihuN/2CXUlH5CJMBFY4ujzi8zZTnAUQZxPetOtDtLyo4prwoWTy', 'killer@mail.ru'); --pass = 3223

CREATE TYPE cyrrency_code AS ENUM ('RUB', 'USD', 'EUR');

CREATE TABLE wallet (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id TEXT NOT NULL, 
    currency cyrrency_code NOT NULL, 
    count DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (user_id) REFERENCES account(uuid) ON DELETE NO ACTION
);

CREATE UNIQUE INDEX idx_wallet_cols ON wallet (user_id, currency);
ALTER TABLE wallet ADD CONSTRAINT uniq_wallet_user_cur UNIQUE (user_id, currency);
ALTER TABLE wallet ADD CONSTRAINT wallet_nonnegative CHECK (count >= 0);

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
    