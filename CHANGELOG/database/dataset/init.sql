CREATE TABLE datasets (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    raw_data_type VARCHAR(255),
    annotation_template_type VARCHAR(255),
    annotation_template_id VARCHAR(255),
    association_id VARCHAR(255) NOT NULL,
    cover_image_url VARCHAR(255),
    UNIQUE KEY association_name (association_id, name)
);
CREATE INDEX datasets_created_at ON datasets (created_at);

CREATE TABLE dataset_versions (
    dataset_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    train_raw_data_view_id VARCHAR(255),
    train_annotation_view_id VARCHAR(255),
    train_raw_data_num INT,
    train_total_data_size BIGINT,
    train_raw_data_ratio FLOAT,
    val_raw_data_view_id VARCHAR(255),
    val_annotation_view_id VARCHAR(255),
    val_raw_data_num INT,
    val_total_data_size BIGINT,
    val_raw_data_ratio FLOAT,
    test_raw_data_view_id VARCHAR(255),
    test_annotation_view_id VARCHAR(255),
    test_raw_data_num INT,
    test_total_data_size BIGINT,
    test_raw_data_ratio FLOAT,
    PRIMARY KEY (dataset_id, name)
);

CREATE TABLE data_pools (
    dataset_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    raw_data_view_id VARCHAR(255),
    annotation_view_id VARCHAR(255),
    PRIMARY KEY (dataset_id, name)
);