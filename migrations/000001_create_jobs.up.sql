CREATE TABLE jobs (
    id CHAR(36) NOT NULL,

    type VARCHAR(100) NOT NULL,

    payload JSON NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'pending',

    attempts INT NOT NULL DEFAULT 0,

    max_attempts INT NOT NULL DEFAULT 3,

    lease_owner VARCHAR(100) NULL,

    lease_expires_at DATETIME(6) NULL,

    last_heartbeat_at DATETIME(6) NULL,

    available_at DATETIME(6) NOT NULL,

    last_error TEXT NULL,

    created_at DATETIME(6) NOT NULL,

    updated_at DATETIME(6) NOT NULL,

    PRIMARY KEY (id),

    INDEX idx_jobs_available (
        status,
        available_at
    ),

    INDEX idx_jobs_lease (
        status,
        lease_expires_at
    )
);