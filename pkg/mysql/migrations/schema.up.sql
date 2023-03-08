CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name text NOT NULL,
    password text NOT NULL,
    designation text NOT NULL,
    employee_id text NOT NULL,
    created_at TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);