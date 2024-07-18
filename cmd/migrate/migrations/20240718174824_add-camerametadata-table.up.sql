CREATE TABLE IF NOT EXISTS camera_metadata (
    cam_id               VARCHAR(255) PRIMARY KEY,
    image_id             VARCHAR(255) NOT NULL,
    camera_name          VARCHAR(255) NOT NULL,
    firmware_version     VARCHAR(255),
    container_name       VARCHAR(255),
    name_of_stored_picture VARCHAR(255),
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL,
    onboarded_at         TIMESTAMP WITH TIME ZONE,
    initialized_at       TIMESTAMP WITH TIME ZONE
);
