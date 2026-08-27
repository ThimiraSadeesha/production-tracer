DROP PROCEDURE IF EXISTS shift_update;
CREATE PROCEDURE shift_update(
    IN shift_id_val BIGINT,
    IN shift_name_val VARCHAR(255),
    IN start_time_val DATETIME(3),
    IN end_time_val DATETIME(3),
    IN department_id_val BIGINT,
    IN status_val VARCHAR(100),
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


    UPDATE tbl_shift
    SET shift = shift_name_val,
        start_time = start_time_val,
        end_time = end_time_val,
        department_id = department_id_val,
        status = status_val,
        updated_by = updated_by_val
    WHERE id = shift_id_val;

    COMMIT;
END;
