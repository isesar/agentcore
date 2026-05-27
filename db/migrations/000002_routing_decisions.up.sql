CREATE TABLE IF NOT EXISTS routing_decisions (
    id SERIAL PRIMARY KEY,
    conversation_id VARCHAR(128) NOT NULL,
    agent_run_id INTEGER REFERENCES agent_runs(id) ON DELETE SET NULL,
    intent VARCHAR(64) NOT NULL,
    confidence DOUBLE PRECISION NOT NULL,
    reason TEXT NOT NULL,
    classifier VARCHAR(32) NOT NULL,
    policy_version VARCHAR(32) NOT NULL,
    trace_id VARCHAR(128) NOT NULL DEFAULT '',
    user_id VARCHAR(128) NOT NULL DEFAULT '',
    source VARCHAR(128) NOT NULL DEFAULT '',
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_routing_decisions_created_at ON routing_decisions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_routing_decisions_intent_created_at ON routing_decisions(intent, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_routing_decisions_trace_created_at ON routing_decisions(trace_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_routing_decisions_conversation_created_at ON routing_decisions(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_routing_decisions_tags_gin ON routing_decisions USING GIN(tags);
