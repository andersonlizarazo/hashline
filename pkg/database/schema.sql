CREATE TABLE token_pool (
    -- Make this table a singleton.
    id INTEGER NOT NULL PRIMARY KEY CHECK (id = 0),
    ring_cursor INTEGER NOT NULL CHECK (ring_cursor >= 0),
    emitted_at INTEGER NOT NULL
) STRICT;

INSERT INTO token_pool (id, ring_cursor, emitted_at) VALUES (0, 0, unixepoch());

CREATE TABLE metadata (
    path TEXT NOT NULL PRIMARY KEY,
    fnv64_hash BLOB NOT NULL,  -- FNV-1 64 over the file bytes, content only.
    updated_at INTEGER NOT NULL   -- UNIX timestamp.
) STRICT;

CREATE TABLE line (
    path TEXT NOT NULL REFERENCES metadata (path) ON DELETE CASCADE,
    token TEXT NOT NULL,
    pos INTEGER NOT NULL,  -- Current line index.
    PRIMARY KEY (path, token),
    CHECK (length(token) BETWEEN 2 AND 8)
) STRICT;
