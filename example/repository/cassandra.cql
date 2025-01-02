CREATE KEYSPACE test WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
USE test;
CREATE TABLE users(
    id UUID,
    name varchar,
    email varchar,
    PRIMARY KEY(id)
);
CREATE INDEX ON users(name);
