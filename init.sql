CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS wallet (
    uuid   uuid PRIMARY KEY,
    amount numeric(20,2) NOT NULL DEFAULT 0
);

INSERT INTO wallet (uuid, amount) VALUES ('31cd8feb-4bd6-424f-b3cb-af52ed07f7dc', 1000.00)
ON CONFLICT (uuid) DO NOTHING;
