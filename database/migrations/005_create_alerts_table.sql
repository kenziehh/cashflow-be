CREATE TABLE alerts (
    id CHAR(26) PRIMARY KEY,
    user_id UUID NOT NULL,
    message VARCHAR(255),
    type VARCHAR(20),
    triggered_at TIMESTAMP,
    CONSTRAINT fk_alerts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);