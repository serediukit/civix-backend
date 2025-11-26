-- Увімкнути розширення
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS postgis;


--------------------------------------------------------
-- TABLE: cities
--------------------------------------------------------
CREATE TABLE IF NOT EXISTS cities
(
    city_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name     VARCHAR(255) NOT NULL,
    region   VARCHAR(255),
    location geometry(Point, 4326)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_cities_location_gist ON cities USING GIST (location);


--------------------------------------------------------
-- TABLE: users
--------------------------------------------------------
CREATE TABLE IF NOT EXISTS users
(
    user_id       SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    name          VARCHAR(100)        NOT NULL,
    surname       VARCHAR(100),
    phone_number  VARCHAR(50) UNIQUE,
    avatar_url    VARCHAR(255),
    reg_city_id   UUID REFERENCES cities (city_id),
    reg_time      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    upd_time      TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    del_time      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_reg_city_id ON users (reg_city_id);


--------------------------------------------------------
-- TABLE: reports
--------------------------------------------------------
CREATE TABLE IF NOT EXISTS reports
(
    report_id         UUID PRIMARY KEY               DEFAULT gen_random_uuid(),
    user_id           INT                   REFERENCES users (user_id) ON DELETE SET NULL,
    create_time       TIMESTAMPTZ                    DEFAULT NOW(),
    update_time       TIMESTAMPTZ                    DEFAULT NOW(),
    location          geometry(Point, 4326) NOT NULL,
    city_id           UUID REFERENCES cities (city_id),
    description       TEXT,
    category_id       INT                   NOT NULL DEFAULT 0,
    current_status_id INT                   NOT NULL DEFAULT 0,
    photo_url         text
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_reports_location_gist ON reports USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports (user_id);
CREATE INDEX IF NOT EXISTS idx_reports_category_id ON reports (category_id);
CREATE INDEX IF NOT EXISTS idx_reports_current_status ON reports (current_status_id);


--------------------------------------------------------
-- TABLE: reports_statuses_log
--------------------------------------------------------
CREATE TABLE IF NOT EXISTS reports_statuses_log
(
    log_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id   UUID NOT NULL REFERENCES reports (report_id) ON DELETE CASCADE,
    status_id   INT,
    update_time TIMESTAMPTZ      DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_statuses_log_report_id ON reports_statuses_log (report_id);
CREATE INDEX IF NOT EXISTS idx_statuses_log_status_id ON reports_statuses_log (status_id);


--------------------------------------------------------
-- TABLE: sessions
--------------------------------------------------------
CREATE TABLE IF NOT EXISTS sessions
(
    session_uid      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          INT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    session_time     TIMESTAMPTZ      DEFAULT NOW(),
    duration_seconds INT,
    city_id          UUID REFERENCES cities (city_id),
    location         geometry(Point, 4326)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_city_id ON sessions (city_id);


--------------------------------------------------------
-- TABLE: ui_actions
--------------------------------------------------------
CREATE TABLE IF NOT EXISTS ui_actions
(
    action_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     INT  NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    session_uid UUID NOT NULL REFERENCES sessions (session_uid) ON DELETE CASCADE,
    action_dt   TIMESTAMPTZ      DEFAULT NOW(),
    context     VARCHAR(255),
    target_info VARCHAR(255)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_actions_user_id ON ui_actions (user_id);
CREATE INDEX IF NOT EXISTS idx_actions_session_uid ON ui_actions (session_uid);
