-- create db if it doesn't already exist
SELECT 'CREATE DATABASE certisafe'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'certisafe');
DROP TABLE IF EXISTS certificates;
CREATE TABLE certificates(
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(30)
);
INSERT INTO certificates(name) VALUES('asd');