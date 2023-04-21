-- create db if it doesn't already exist
SELECT 'CREATE DATABASE certifisafe'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'certifisafe');
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS passwordRecovery;
DROP TABLE IF EXISTS verifications;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    password VARCHAR(100),
    first_name VARCHAR(30),
    last_name VARCHAR(30),
    phone VARCHAR(30),
    is_admin BOOLEAN,
    is_active BOOLEAN
);

CREATE TABLE passwordRecovery (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    code VARCHAR(100),
    is_used BOOLEAN
);

CREATE TABLE verifications (
      id SERIAL PRIMARY KEY,
      email VARCHAR(50),
      code VARCHAR(100)
);

CREATE TABLE certificates (
     id  BIGSERIAL PRIMARY KEY,
     name VARCHAR(30),
     valid_from DATE,
     valid_to DATE,
     subject_id INTEGER REFERENCES users(id),
     issuer_id INTEGER REFERENCES users(id),
     status SMALLINT,
     type SMALLINT,
     organization VARCHAR(100)
);

CREATE TABLE requests (
      id SERIAL PRIMARY KEY,
      datetime DATE,
      parent_certificate_id SMALLINT REFERENCES certificates(id),
      certificate_id SMALLINT NOT NULL REFERENCES certificates(id),
      status SMALLINT
);

INSERT INTO users(email, password, first_name, last_name, phone, is_admin)
VALUES('project.usertest+sladic@outlook.com', '$2a$12$u9LD12t.4WxM/nmMiNCB2e0Tj9pVfQcSyJiIzm4vMvEl/zemkKoee', 'Goran', 'Sladic', '065482564', true);


