CREATE TABLE actions (
    id VARCHAR(255),
    name VARCHAR(255),
    resource_management_id VARCHAR(255),
    resource_type VARCHAR(255),
    resource_id VARCHAR(255),
    action VARCHAR(64),
    params VARCHAR(2048),
    create_at BIGINT,
    update_at BIGINT,
    PRIMARY KEY (resource_type, resource_id)
);
