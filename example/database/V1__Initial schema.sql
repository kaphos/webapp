CREATE TABLE users (
    id     SERIAL PRIMARY KEY,
    kc_sub TEXT NOT NULL
);

CREATE TABLE items (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    email TEXT NOT NULL
);

INSERT INTO items (name, email)
    VALUES ('John', 'john@gmail.com'),
           ('Tom', 'tom@email.com'),
           ('Jane', 'jane@hotmail.com');