CREATE TABLE IF NOT EXISTS bookmarks (
    id          UUID         PRIMARY KEY,
    url         TEXT         NOT NULL,
    title       TEXT         NOT NULL,
    tags        TEXT[]       NOT NULL DEFAULT '{}',
    notes       TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks (created_at DESC);