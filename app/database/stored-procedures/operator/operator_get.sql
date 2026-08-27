DROP PROCEDURE IF EXISTS operator_get;
CREATE PROCEDURE operator_get(
    IN operator_id_val BIGINT
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

    SELECT o.id              AS operatorId,
           o.emp_no          AS empNo,
           o.name            AS operatorName,
           o.section         AS section,
           o.status          AS operatorStatus,
           d.department_name AS departmentName,
           d.id              AS departmentId
    FROM tbl_operator o
             JOIN tbl_department d ON o.department_id = d.id
    WHERE o.id = operator_id_val;
    COMMIT;
END;