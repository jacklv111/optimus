CREATE TABLE resource_management (
    id VARCHAR(255) PRIMARY KEY,
    domain VARCHAR(255),
    workspace varchar(255),
    created_at BIGINT,
    CONSTRAINT domain_workspace UNIQUE (domain, workspace)
);
