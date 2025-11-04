CREATE TABLE IF NOT EXISTS metric.counters (
    id SERIAL PRIMARY KEY,
    metric_type TEXT NOT NULL DEFAULT 'counter',
    metric_name TEXT UNIQUE NOT NULL,
    metric_value BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS metric.gauges (
    id SERIAL PRIMARY KEY,
    metric_type TEXT NOT NULL DEFAULT 'gauge',
    metric_name TEXT UNIQUE NOT NULL,
    metric_value DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);



