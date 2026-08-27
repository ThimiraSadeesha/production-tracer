DROP PROCEDURE IF EXISTS operator_update;
CREATE PROCEDURE operator_update(
    IN operator_id_val BIGINT,
    IN emp_no_val VARCHAR(255),
    IN name_val VARCHAR(255),
    IN section_val VARCHAR(255),
    IN status_val VARCHAR(100),
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

    UPDATE tbl_operator
    SET emp_no        = emp_no_val,
        name          = name_val,
        section       = section_val,
        status        = status_val,
        department_id = department_id_val,
        updated_by    = updated_by_val
    WHERE id = operator_id_val;
    COMMIT;

    SELECT operator_id_val AS operatorId;
END;
