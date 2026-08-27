DROP PROCEDURE IF EXISTS department_find;
CREATE PROCEDURE department_find(
    IN department_code_val VARCHAR(50),
    IN department_name_val VARCHAR(255),
    IN department_status_val VARCHAR(255),
    IN p_cursor BIGINT,
    IN p_limit INT
)
BEGIN

    IF p_limit IS NULL OR p_limit <= 0 THEN
        SET p_limit = 20;
    END IF;

    SELECT id,
           department_code As departmentCode,
           department_name As departmentName,
           department_status  status
    FROM tbl_department
    WHERE (department_code_val IS NULL OR department_code_val = '' OR
           department_code LIKE CONCAT('%', department_code_val, '%'))
      AND (department_name_val IS NULL OR department_name_val = '' OR
           department_name LIKE CONCAT('%', department_name_val, '%'))
      AND (department_status_val IS NULL OR department_status_val = '' OR
           department_status LIKE CONCAT('%', department_status_val, '%'))
      AND (p_cursor = -1 OR id < p_cursor)
    ORDER BY id DESC
    LIMIT p_limit;
END;