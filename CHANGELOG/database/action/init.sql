CREATE TABLE actions (
    name VARCHAR(255),
    resource_management_id VARCHAR(255),
    resource_type VARCHAR(255),
    resource_id VARCHAR(255),
    params VARCHAR(2048),
    create_at BIGINT,
    update_at BIGINT,
    PRIMARY KEY (resource_type, resource_id)
);
