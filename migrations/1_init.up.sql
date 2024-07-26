
CREATE TABLE IF NOT EXISTS message
(
    id SERIAL PRIMARY KEY,
    content VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS outbox
(
    id SERIAL PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    status VARCHAR(4) CHECK(status in ('new', 'done')) NOT NULL,
    cerate_at TIMESTAMP NOT NULL,
    reserved TIMESTAMP NOT NULL
);