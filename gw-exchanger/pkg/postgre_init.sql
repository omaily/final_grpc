SET client_encoding = 'UTF8';

CREATE TABLE account (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid TEXT NOT NULL,
    name TEXT NOT NULL, 
    amount INT
);

CREATE UNIQUE INDEX idx_account_uuid ON account(uuid);

INSERT INTO account (uuid, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000001', 'Алеша', 52);
    
INSERT INTO account (uuid, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000002', 'Степаша', 322);

INSERT INTO account (uuid, name, amount) 
    VALUES ('66666666-6666-6666-6666-666666666666', 'Олег', 666);

INSERT INTO account (uuid, name, amount) 
    VALUES ('00000000-0000-0000-0000-000000000003', 'Никита', 3141592);