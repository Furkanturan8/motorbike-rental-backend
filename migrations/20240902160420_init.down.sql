-- Add down migration script here

-- Drop Auth Table
DROP TABLE IF EXISTS auth_refresh_tokens;

-- Drop Users Table
DROP TABLE IF EXISTS users;

-- Drop Motorbike Table
DROP TABLE IF EXISTS motorbike;

DROP TABLE IF EXISTS motorbike_photos;