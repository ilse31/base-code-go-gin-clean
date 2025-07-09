-- Drop indexes
DROP INDEX IF EXISTS httplog.idx_log_outgoing_requests_trace_id;
DROP INDEX IF EXISTS httplog.idx_log_outgoing_requests_created_at;
DROP INDEX IF EXISTS httplog.idx_log_incoming_requests_trace_id;
DROP INDEX IF EXISTS httplog.idx_log_incoming_requests_created_at;
DROP INDEX IF EXISTS httplog.idx_log_errors_trace_id;
DROP INDEX IF EXISTS httplog.idx_log_errors_created_at;

-- Drop tables in reverse order of creation to respect foreign key constraints
DROP TABLE IF EXISTS httplog.log_errors;
DROP TABLE IF EXISTS httplog.log_incoming_requests;
DROP TABLE IF EXISTS httplog.log_outgoing_requests;

-- Drop the http_method type if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'http_method') THEN
        DROP TYPE http_method;
    END IF;
END$$;

-- Drop the schema if it's empty
DROP SCHEMA IF EXISTS httplog;
