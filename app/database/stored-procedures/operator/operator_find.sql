DROP PROCEDURE IF EXISTS operator_find;
CREATE PROCEDURE operator_find(
    IN emp_no_val VARCHAR(255),
    IN name_val VARCHAR(255),
    IN section_val VARCHAR(255),
    IN status_val VARCHAR(100),
    IN department_id_val BIGINT,
    IN p_cursor BIGINT,
    IN p_limit INT
)
BEGIN

    IF p_limit IS NULL OR p_limit <= 0 THEN
        SET p_limit = 20;
    END IF;

    SELECT o.id               AS operatorId,
           o.emp_no           AS empNo,
           o.name             AS operatorName,
           o.section          AS section,
           o.status           AS operatorStatus,
           td.department_name AS departmentName
    FROM tbl_operator o
             JOIN tbl_department td ON o.department_id = td.id
    WHERE (emp_no_val IS NULL OR emp_no_val = '' OR o.emp_no LIKE CONCAT('%', emp_no_val, '%'))
      AND (name_val IS NULL OR name_val = '' OR o.name LIKE CONCAT('%', name_val, '%'))
      AND (section_val IS NULL OR section_val = '' OR o.section LIKE CONCAT('%', section_val, '%'))
      AND (status_val IS NULL OR status_val = '' OR o.status LIKE CONCAT('%', status_val, '%'))
      AND (department_id_val IS NULL OR department_id_val = -1 OR o.department_id = department_id_val)
      AND (p_cursor = -1 OR o.id < p_cursor)
    ORDER BY o.id DESC
    LIMIT p_limit;
END;
