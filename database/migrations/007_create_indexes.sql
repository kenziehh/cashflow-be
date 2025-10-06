CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_alerts_user ON alerts(user_id);
CREATE INDEX idx_maximum_spends_user ON maximum_spends(user_id);