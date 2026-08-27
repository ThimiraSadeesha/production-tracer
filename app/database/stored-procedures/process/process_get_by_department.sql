DROP PROCEDURE IF EXISTS process_get_by_department;
CREATE PROCEDURE process_get_by_department(
    IN department_id_val BIGINT
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

    SELECT td.id                AS departmentId,
           td.department_code   AS departmentCode,
           td.department_name   AS departmentName,
           td.department_status AS departmentStatus,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'processId', tp.id,
                           'processCode', tp.code,
                           'processName', tp.name,
                           'processStatus', tp.status,
                           'sequence', tp.sequence,
                           'isRequired', tp.is_required
                   )
           )                    AS processes
    FROM tbl_department td
             LEFT JOIN tbl_process tp ON tp.department_id = td.id
    WHERE td.id = department_id_val
      AND tp.status = 'Active'
    GROUP BY td.id;
    COMMIT;
END;
