CREATE TABLE documents (
    id CHAR(36) NOT NULL,

    filename VARCHAR(255) NOT NULL,
    content_type VAR_CHAT(255) NOT NULL,
    size_bytes BIGINY UNSIGNED NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'pending',

    last_error TEXT NULL,

    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    PRIMARY KEY (id),

    INDEX idx_documents_status (status),
    INDEX idx_documents_created_at (created_at0)
);