DROP PROCEDURE IF EXISTS operator_save;
CREATE PROCEDURE operator_save(
    IN emp_no_val VARCHAR(255),
    IN name_val VARCHAR(255),
    IN section_val VARCHAR(255),
    IN status_val VARCHAR(100),
    IN department_id_val BIGINT,
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

    INSERT INTO tbl_operator (emp_no, name, section, status, department_id, created_by)
    VALUES (emp_no_val, name_val, section_val, status_val, department_id_val, created_by_val);

    COMMIT;

    SELECT LAST_INSERT_ID() AS operatorId;
END;