-- Add up migration script here

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    surname VARCHAR(20) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role INT NOT NULL CHECK (role IN (1, 10)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Auth_refresh_tokens Table
CREATE TABLE IF NOT EXISTS auth_refresh_tokens (
    token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id BIGINT NOT NULL,
    role BIGINT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS motorbike (
    id SERIAL PRIMARY KEY,
    model VARCHAR(100) NOT NULL,
    location_latitude DOUBLE PRECISION NOT NULL,
    location_longitude DOUBLE PRECISION NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('available', 'maintenance', 'rented')),
    lock_status VARCHAR(10) NOT NULL CHECK (lock_status IN ('locked', 'unlocked')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS motorbike_photos (
   id SERIAL PRIMARY KEY,
   motorbike_id INTEGER NOT NULL REFERENCES motorbike(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS rides (
     id SERIAL PRIMARY KEY,
     user_id INT NOT NULL,
     motorbike_id INT NOT NULL,
     start_time TIMESTAMPTZ NOT NULL,
     end_time TIMESTAMPTZ,
     duration INTERVAL NOT NULL,
     cost NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (motorbike_id) REFERENCES motorbike(id) ON DELETE RESTRICT
);

-- Optional: Create indexes on foreign keys for better performance
CREATE INDEX idx_rides_user_id ON rides(user_id);
CREATE INDEX idx_rides_motorbike_id ON rides(motorbike_id);


-- Insert default admin user
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        INSERT INTO users (name, surname, username, phone, email, password, role, created_at, updated_at)
        VALUES (
            'Furkan',
            'Turan',
            'admin',
            '1234567890',
            'admin@example.com',
            '$2a$14$4aHVjRPGxCSpvNGM7tm6COHMDZ5LRzk/ehW0A6AOxoEUcnyYQgbue',
            10,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        );
END IF;
END
$$;

