CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS camera_metadata (
    cam_id               VARCHAR(36) NOT NULL DEFAULT uuid_generate_v4()::text,
    image_id             VARCHAR(36) ,
    camera_name          VARCHAR(255) NOT NULL,
    firmware_version     VARCHAR(255) NOT NULL,
    container_name       VARCHAR(255),
    name_of_stored_picture VARCHAR(255),
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL,
    onboarded_at         TIMESTAMP WITH TIME ZONE,
    initialized_at       TIMESTAMP WITH TIME ZONE
);
