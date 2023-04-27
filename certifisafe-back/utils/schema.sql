-- create db if it doesn't already exist
SELECT 'CREATE DATABASE certifisafe'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'certifisafe');
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS passwordRecovery;


CREATE TABLE passwordRecovery (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    code VARCHAR(100),
    is_used BOOLEAN
);


--INSERT INTO users(email, password, first_name, last_name, phone, is_admin)
--VALUES('project.usertest+sladic@outlook.com', '$2a$12$u9LD12t.4WxM/nmMiNCB2e0Tj9pVfQcSyJiIzm4vMvEl/zemkKoee', 'Goran', 'Sladic', '065482564', true);


