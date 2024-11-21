CREATE TABLE "users" (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
);

CREATE TABLE "events" (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id INTEGER NOT NULL,
    kind TEXT NOT NULL,
    version TEXT NOT NULL,
    created_at timestamp with time zone default current_timestamp,
    data JSONB NOT NULL
);