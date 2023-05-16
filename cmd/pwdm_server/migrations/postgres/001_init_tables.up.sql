CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(uuid UUID UNIQUE NOT NULL PRIMARY KEY, login VARCHAR(255) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL, deleted BOOLEAN DEFAULT false);
CREATE TABLE IF NOT EXISTS log_pwd_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), login VARCHAR(255), password VARCHAR(255), tag VARCHAR(255), comment TEXT, deleted BOOLEAN DEFAULT false);
CREATE TABLE IF NOT EXISTS card_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), num VARCHAR(255), date VARCHAR(255), cvc VARCHAR(255), first_name VARCHAR(255), last_name VARCHAR(255), tag VARCHAR(255), comment TEXT, deleted BOOLEAN DEFAULT false);
CREATE TABLE IF NOT EXISTS text_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), data TEXT, tag VARCHAR(255), comment TEXT, deleted BOOLEAN DEFAULT false);
CREATE TABLE IF NOT EXISTS binary_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), data TEXT, tag VARCHAR(255), comment TEXT, deleted BOOLEAN DEFAULT false);
