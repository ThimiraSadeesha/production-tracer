DROP PROCEDURE IF EXISTS process_update;
CREATE PROCEDURE process_update(
    IN process_id_val BIGINT,
    IN code_val VARCHAR(100),
    IN name_val VARCHAR(255),
    IN status_val VARCHAR(255),
    IN is_required_val TINYINT(1),
    IN department_id_val BIGINT,
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


    SELECT COUNT(*)
    INTO @record_count
    FROM tbl_process
    WHERE code = code_val
      AND id != process_id_val;

    IF @record_count > 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Process code already exists for another record.';
    END IF;

    UPDATE tbl_process
    SET code          = code_val,
        name          = name_val,
        status        = status_val,
        is_required   = is_required_val,
        department_id = department_id_val,
        updated_by    = updated_by_val
    WHERE id = process_id_val;
    COMMIT;
END;