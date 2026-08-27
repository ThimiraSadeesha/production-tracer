DROP PROCEDURE IF EXISTS shift_get;
CREATE PROCEDURE shift_get(
    IN shift_id_val BIGINT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            ROLLBACK;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
        END;

    START TRANSACTION;

    SELECT ts.id            AS shiftId,
           ts.shift         AS shiftName,
           ts.start_date    AS startDate,
           ts.end_date      AS endDate,
           ts.start_time    AS startTime,
           ts.end_time      AS endTime,
           ts.department_id AS departmentId,
           ts.status        AS shiftStatus
    FROM tbl_shift ts
    WHERE ts.id = shift_id_val;
    COMMIT;
END;
