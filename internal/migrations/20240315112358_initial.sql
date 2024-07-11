-- +goose Up
-- +goose StatementBegin
SET TIMEZONE = 'Europe/Moscow';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    is_superuser BOOLEAN DEFAULT FALSE,
    is_staff BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS sessions(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES users ON DELETE CASCADE,
    refresh_token VARCHAR(1000) NOT NULL
);

CREATE TABLE IF NOT EXISTS portfolios(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES users ON DELETE CASCADE,
    cmc_cryptocurrency_id INT NOT NULL,
    cryptocurrency VARCHAR(255),
    cryptocurrency_symbol VARCHAR(30),
    price FLOAT NOT NULL,
    count FLOAT NOT NULL,
    purchase_time TIMESTAMP NOT NULL,
    commentary VARCHAR,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS portfolios CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
