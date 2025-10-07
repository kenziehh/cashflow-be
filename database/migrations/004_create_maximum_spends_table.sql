CREATE TABLE maximum_spends (
    id CHAR(26) PRIMARY KEY,
    user_id UUID NOT NULL,
    daily_limit DECIMAL(12,2),
    monthly_limit DECIMAL(12,2),
    yearly_limit DECIMAL(12,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_maximum_spends_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
