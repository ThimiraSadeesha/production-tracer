DROP PROCEDURE IF EXISTS department_save;
CREATE PROCEDURE department_save(
    IN dep_code VARCHAR(255),
    IN dep_name VARCHAR(255),
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE record_count INT DEFAULT 0;
    DECLARE inserted_id INT;

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
    INTO record_count
    FROM tbl_department
    WHERE department_code = dep_code;

    IF record_count > 0 THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'Department-Code already exists.';
    END IF;
    INSERT INTO tbl_department (department_code,
                                department_name,
                                department_status,
                                created_by,
                                created_at)
    VALUES (dep_code,
            dep_name,
            'Active',
            created_by_val,
            NOW());
    SET inserted_id = LAST_INSERT_ID();
    COMMIT;
    SELECT inserted_id AS departmentId;
END;
