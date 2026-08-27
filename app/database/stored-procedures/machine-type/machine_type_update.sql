DROP PROCEDURE IF EXISTS machine_type_update;
CREATE PROCEDURE machine_type_update(
    IN mt_id BIGINT,
    IN mt_code VARCHAR(100),
    IN mt_name VARCHAR(255),
    IN mt_status VARCHAR(255),
    IN mt_process_id BIGINT,
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

    UPDATE tbl_machine_type
    SET code                = mt_code,
        name                = mt_name,
        machine_type_status = mt_status,
        process_id          = mt_process_id,
        updated_by          = updated_by_val
    WHERE id = mt_id;
    COMMIT;

    SELECT mt_id AS updated_machine_type_id;
END;