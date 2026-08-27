DROP PROCEDURE IF EXISTS shift_save;
CREATE PROCEDURE shift_save(
    IN shift_name_val VARCHAR(255),
    IN start_date_val DATE,
    IN start_time_val TIME,
    IN end_date_val DATE,
    IN end_time_val TIME,
    IN department_id_val BIGINT,
    IN status_val VARCHAR(100),
    IN created_by_val VARCHAR(255)
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

    INSERT INTO tbl_shift (shift,
                           start_date,
                           start_time,
                           end_date,
                           end_time,
                           department_id,
                           status,
                           created_by)
    VALUES (shift_name_val,
            start_date_val,
            start_time_val,
            end_date_val,
            end_time_val,
            department_id_val,
            status_val,
            created_by_val);
    COMMIT;

    SELECT LAST_INSERT_ID() AS shiftId;
END;