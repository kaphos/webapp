CREATE TABLE users (
    id     SERIAL PRIMARY KEY,
    name   TEXT  NOT NULL,
    email  TEXT  NOT NULL,
    admin  BOOL  NOT NULL,
    groups INT   NOT NULL,
    age    FLOAT NOT NULL
);

CREATE TABLE items (
    id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    edited  TIMESTAMP,
    name    TEXT      NOT NULL,
    owner   TEXT,
    found   BOOL,
    count   INT,
    price   FLOAT
);

INSERT INTO users (name, email, admin, groups, age)
    VALUES ('John', 'john@gmail.com', TRUE, 1, 3.5),
           ('Tom', 'tom@email.com', FALSE, 2, 5.2),
           ('Jane', 'jane@hotmail.com', FALSE, 3, 8.9);

INSERT INTO items (name, owner, found, count, price)
    VALUES ('Item Name', 'TheOwner', TRUE, 0, 10.0),
           ('Second Item', NULL, NULL, NULL, NULL);