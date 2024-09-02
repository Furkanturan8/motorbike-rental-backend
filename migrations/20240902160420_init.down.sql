-- Add down migration script here

-- Drop Auth Table
DROP TABLE IF EXISTS auth_refresh_tokens;

-- Drop Users Table
DROP TABLE IF EXISTS users;

