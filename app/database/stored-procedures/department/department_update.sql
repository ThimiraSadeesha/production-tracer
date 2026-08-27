DROP PROCEDURE IF EXISTS department_update;
CREATE PROCEDURE department_update(
    IN dep_id INT,
    IN dep_code VARCHAR(255),
    IN dep_name VARCHAR(255),
    IN dep_status VARCHAR(100),
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

    UPDATE tbl_department
    SET department_code = IFNULL(dep_code, department_code),
        department_name = IFNULL(dep_name, department_name),
        department_status = IFNULL(dep_status, department_status),
        updated_by = updated_by_val,
        updated_at = NOW()
    WHERE id = dep_id;

    COMMIT;
    SELECT dep_id AS departmentId;
END;
