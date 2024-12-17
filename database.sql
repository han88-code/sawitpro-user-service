CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	fullname VARCHAR(60) NOT NULL,
	phonenumber VARCHAR(13) NOT NULL,
	password VARCHAR(256) NOT NULL,
	successlogin INTEGER DEFAULT 0
);
INSERT INTO users(fullname,phonenumber,password) VALUES ('fullname01','phonenumber01','password01');