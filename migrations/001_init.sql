-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscriptions (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name VARCHAR(255)     NOT NULL,
    price        INTEGER          NOT NULL CHECK (price > 0),
    user_id      UUID             NOT NULL,
    start_date   DATE             NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id      ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date   ON subscriptions(start_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date     ON subscriptions(end_date);

-- Дополнительная проверка (опционально, но полезно)
ALTER TABLE subscriptions 
ADD CONSTRAINT chk_end_date_after_start 
CHECK (end_date IS NULL OR end_date >= start_date);

-- +goose Down
DROP TABLE IF EXISTS subscriptions;