CREATE DATABASE IF NOT EXISTS user_auth;

USE user_auth;

CREATE TABLE IF NOT EXISTS users(
   id SERIAL,
   email VARCHAR(50) NOT NULL,
   password CHAR(60) NOT NULL,
   name VARCHAR(50),
   address VARCHAR(100),
   phone VARCHAR(30),
   token CHAR(100),
   recovery_token CHAR(100),
   google_id VARCHAR(50)
);