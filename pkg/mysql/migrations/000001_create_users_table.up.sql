CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name text NOT NULL,
    password text NOT NULL,
    designation text NOT NULL,
    employee_id text NOT NULL,
    is_admin int NOT NULL,
    is_active int NOT NULL,
    created_at TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reservations (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    reservation_time TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    type BIGINT NOT NULL,
    no_of_guests BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);