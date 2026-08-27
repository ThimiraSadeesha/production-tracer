DROP PROCEDURE IF EXISTS department_get_all;
CREATE PROCEDURE department_get_all()
BEGIN
    SELECT d.id                AS id,
           d.department_code   AS departmentCode,
           d.department_name   AS departmentName,
           d.department_status AS status
    FROM tbl_department d
    GROUP BY d.id, d.department_code, d.department_name;

END
