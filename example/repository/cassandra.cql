CREATE KEYSPACE roster WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
USE roster;
CREATE TABLE users(
    id UUID,
    name varchar,
    email varchar,
    PRIMARY KEY(id)
);
CREATE INDEX ON users(name);
