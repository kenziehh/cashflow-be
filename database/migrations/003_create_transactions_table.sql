CREATE TYPE transaction_type AS ENUM ('income', 'expense');

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    category_id ULID,
    type transaction_type NOT NULL,
    period VARCHAR(20),
    amount DECIMAL(12,2) NOT NULL,
    note TEXT,
    date DATE NOT NULL,
    proof_file VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transactions_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);