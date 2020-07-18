-- CREATE USER "user_auth";
-- ALTER USER "user_auth" WITH ENCRYPTED PASSWORD 'user_authpass';

CREATE DATABASE IF NOT EXISTS user_auth;
-- GRANT ALL PRIVILEGES ON DATABASE "user_auth" TO "user_auth";
-- \c "user_auth"
-- CREATE EXTENSION pgcrypto;