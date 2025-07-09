-- Create schema for HTTP logging
CREATE SCHEMA IF NOT EXISTS httplog;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create enum for HTTP methods if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'http_method') THEN
        CREATE TYPE http_method AS ENUM (
            'GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS', 'TRACE', 'CONNECT'
        );

END IF;

END$$;

-- Create log_outgoing_requests table
CREATE TABLE IF NOT EXISTS httplog.log_outgoing_requests (
    id SERIAL PRIMARY KEY,
    trace_id TEXT NOT NULL,
    event_name TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    METHOD http_method NOT NULL,
    request JSONB NOT NULL,
    status_code INTEGER NOT NULL,
    response JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create log_incoming_requests table
CREATE TABLE IF NOT EXISTS httplog.log_incoming_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trace_id TEXT NOT NULL,
    event_name TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    METHOD http_method NOT NULL,
    request JSONB NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create log_errors table
CREATE TABLE IF NOT EXISTS httplog.log_errors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trace_id TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    error TEXT NOT NULL,
    request_id UUID REFERENCES httplog.log_incoming_requests (id) ON DELETE SET NULL,
    traces JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_log_outgoing_requests_trace_id ON httplog.log_outgoing_requests (trace_id);

CREATE INDEX IF NOT EXISTS idx_log_outgoing_requests_created_at ON httplog.log_outgoing_requests (created_at);

CREATE INDEX IF NOT EXISTS idx_log_incoming_requests_trace_id ON httplog.log_incoming_requests (trace_id);

CREATE INDEX IF NOT EXISTS idx_log_incoming_requests_created_at ON httplog.log_incoming_requests (created_at);

CREATE INDEX IF NOT EXISTS idx_log_errors_trace_id ON httplog.log_errors (trace_id);

CREATE INDEX IF NOT EXISTS idx_log_errors_created_at ON httplog.log_errors (created_at);

-- Add comments for documentation
COMMENT ON SCHEMA httplog IS 'Schema for storing HTTP request/response logs and errors';

COMMENT ON TABLE httplog.log_outgoing_requests IS 'Stores logs for outgoing HTTP requests made by the application';

COMMENT ON TABLE httplog.log_incoming_requests IS 'Stores logs for incoming HTTP requests to the application';

COMMENT ON TABLE httplog.log_errors IS 'Stores error logs with references to the original requests';

-- Add column comments
COMMENT ON COLUMN httplog.log_outgoing_requests.trace_id IS 'Distributed tracing ID for correlating logs';

COMMENT ON COLUMN httplog.log_incoming_requests.trace_id IS 'Distributed tracing ID for correlating logs';

COMMENT ON COLUMN httplog.log_errors.trace_id IS 'Distributed tracing ID for correlating logs';

COMMENT ON COLUMN httplog.log_errors.request_id IS 'Reference to the original request that caused the error';