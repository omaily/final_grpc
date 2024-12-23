SET client_encoding = 'UTF8';

CREATE TABLE account (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_account_uuid ON account(uuid);

INSERT INTO account (uuid, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'Алеша');
INSERT INTO account (uuid, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000002', 'Степаша');
INSERT INTO account (uuid, name, amount) 
    VALUES ('66666666-6666-6666-6666-666666666666', 'Олег');
INSERT INTO account (uuid, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000003', 'Никита');



CREATE TABLE wallet (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id TEXT NOT NULL, 
    currency TEXT NOT NULL, 
    count DOUBLE PRECISION,
    FOREIGN KEY (user_id) REFERENCES account(uuid) ON DELETE NO ACTION
);

CREATE UNIQUE INDEX idx_wallet_uuid ON wallet(id);

INSERT INTO wallet (user_id, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'EUR', 0);
INSERT INTO wallet (user_id, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'USD', 17.3);
INSERT INTO wallet (user_id, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'RUB', 10000.0);
    