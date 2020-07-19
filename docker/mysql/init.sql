-- CREATE USER "user_auth";
-- ALTER USER "user_auth" WITH ENCRYPTED PASSWORD 'user_authpass';

CREATE DATABASE IF NOT EXISTS user_auth;
-- GRANT ALL PRIVILEGES ON DATABASE "user_auth" TO "user_auth";
-- \c "user_auth"
-- CREATE EXTENSION pgcrypto;
USE user_auth;

CREATE TABLE IF NOT EXISTS users(
   id SERIAL,
   email VARCHAR(50) NOT NULL,
   password CHAR(60) NOT NULL,
   name VARCHAR(50),
   address VARCHAR(100),
   phone VARCHAR(30),
   token CHAR(100)
);