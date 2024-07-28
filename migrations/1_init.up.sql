
CREATE TABLE IF NOT EXISTS message
(
    id SERIAL PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    status VARCHAR(6) CHECK (status IN ('none', 'update')) NOT NULL,
    processed_content VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS outbox
(
    id SERIAL PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    status VARCHAR(4) CHECK (status IN ('new', 'done')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reserved TIMESTAMP
);

CREATE INDEX idx_message_status ON message (status);
CREATE INDEX idx_outbox_status ON outbox (status);