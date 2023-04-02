-- create db if it doesn't already exist
SELECT 'CREATE DATABASE certisafe'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'certisafe');
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    password VARCHAR(100),
    first_name VARCHAR(30),
    last_name VARCHAR(30),
    phone VARCHAR(30),
    is_admin BOOLEAN
);

CREATE TABLE certificates (
     id SERIAL PRIMARY KEY,
     name VARCHAR(30),
     valid_from date,
     valid_to date,
     subject_pk varchar(30),
     subject_id INTEGER REFERENCES users(id),
     issuer_id INTEGER REFERENCES users(id),
     signature  varchar(30)
);


INSERT INTO users(email, password, first_name, last_name, phone, is_admin) VALUES('project.usertest+sladic@outlook.com', '$2a$12$u9LD12t.4WxM/nmMiNCB2e0Tj9pVfQcSyJiIzm4vMvEl/zemkKoee', 'Goran', 'Sladic', '065482564', true);
INSERT INTO users(email, password, first_name, last_name, phone, is_admin) VALUES('project.usertest+majstorovic@outlook.com', '', 'Nemanja', 'Majstorovic', '063622564', false);
INSERT INTO users(email, password, first_name, last_name, phone, is_admin) VALUES('project.usertest+dutina@outlook.com', '', 'Nemanja', 'Dutina', '061882596', false);
INSERT INTO users(email, password, first_name, last_name, phone, is_admin) VALUES('project.usertest+milosavljevic@outlook.com', '', 'Branko', 'Milosavljevic', '0604152368', false);

INSERT INTO certificates(name, valid_from, valid_to, subject_id, subject_pk, issuer_id, signature)
VALUES('Certificate #1', '2022-01-01', '2024-01-01', 2, 'asd', 1, 'asdasd');