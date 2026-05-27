CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role VARCHAR(32) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

CREATE TABLE IF NOT EXISTS agent_runs (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    model VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_agent_runs_conversation_id ON agent_runs(conversation_id);
CREATE INDEX IF NOT EXISTS idx_agent_runs_started_at ON agent_runs(started_at);

CREATE TABLE IF NOT EXISTS tool_calls (
    id SERIAL PRIMARY KEY,
    agent_run_id INTEGER NOT NULL REFERENCES agent_runs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    arguments JSONB NOT NULL DEFAULT '{}'::jsonb,
    result JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tool_calls_agent_run_id ON tool_calls(agent_run_id);

CREATE TABLE IF NOT EXISTS usage_events (
    id SERIAL PRIMARY KEY,
    agent_run_id INTEGER NOT NULL REFERENCES agent_runs(id) ON DELETE CASCADE,
    prompt_tokens INTEGER NOT NULL DEFAULT 0,
    completion_tokens INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_usage_events_agent_run_id ON usage_events(agent_run_id);

CREATE TABLE IF NOT EXISTS errors (
    id SERIAL PRIMARY KEY,
    agent_run_id INTEGER REFERENCES agent_runs(id) ON DELETE SET NULL,
    error_type VARCHAR(128) NOT NULL,
    message TEXT NOT NULL,
    stack_trace TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_errors_agent_run_id ON errors(agent_run_id);
CREATE INDEX IF NOT EXISTS idx_errors_created_at ON errors(created_at);
