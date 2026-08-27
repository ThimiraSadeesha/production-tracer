DROP PROCEDURE IF EXISTS device_update;
CREATE PROCEDURE device_update(
    IN p_device_id BIGINT,
    IN p_app_version VARCHAR(50),
    IN p_last_seen_at LONGTEXT,
    IN p_status VARCHAR(50),
    IN p_device_type VARCHAR(50),
    IN p_machine_id BIGINT,
    IN p_process_id BIGINT,
    IN p_android_id VARCHAR(150),
    IN updated_by_val VARCHAR(255)
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    START TRANSACTION;

    UPDATE tbl_device
    SET app_version  = p_app_version,
        last_seen_at = p_last_seen_at,
        status       = p_status,
        device_type  = COALESCE(NULLIF(p_device_type, ''), device_type),
        machine_id   = p_machine_id,
        process_id   = p_process_id,
        android_id   = p_android_id,
        updated_by   = updated_by_val
    WHERE id = p_device_id;

    SELECT p_device_id AS deviceId;
    COMMIT;
END;