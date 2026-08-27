DROP PROCEDURE IF EXISTS process_get;
CREATE PROCEDURE process_get(
    IN process_id_val INT
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

    SELECT tp.id            AS processId,
           tp.code          AS processCode,
           tp.name          AS processName,
           tp.status        AS processStatus,
           tp.is_required   AS isRequired,
           tp.department_id AS departmentId,
           IFNULL(
                   (SELECT JSON_ARRAYAGG(
                                   JSON_OBJECT(
                                           'id', tm.id,
                                           'code', tm.code,
                                           'name', tm.name,
                                           'status', tm.machine_type_status
                                   )
                           )
                    FROM tbl_machine_type tm
                    WHERE tm.process_id = tp.id),
                   JSON_ARRAY()
           )                AS machineTypes,
           IFNULL(
                   (SELECT JSON_ARRAYAGG(
                                   JSON_OBJECT(
                                           'id', tm.id,
                                           'name', tm.attribute_name,
                                           'isRequired', tm.is_required
                                   )
                           )
                    FROM tbl_process_attributes tm
                    WHERE tm.process_id = tp.id),
                   JSON_ARRAY()
           )                AS attributes
    FROM tbl_process tp
    WHERE tp.id = process_id_val;
    COMMIT;
END;
