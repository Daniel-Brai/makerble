CREATE TABLE IF NOT EXISTS invalid_tokens (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    token TEXT NOT NULL,
    invalidated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_invalid_tokens_token ON invalid_tokens(token);
CREATE INDEX idx_invalid_tokens_expires_at ON invalid_tokens(expires_at);
