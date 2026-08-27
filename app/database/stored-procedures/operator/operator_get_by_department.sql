DROP PROCEDURE IF EXISTS operator_get_by_department;
CREATE PROCEDURE operator_get_by_department(
    IN departmentId BIGINT
)
BEGIN
    SELECT td.id                AS departmentId,
           td.department_code   AS departmentCode,
           td.department_name   AS departmentName,
           td.department_status AS departmentStatus,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'operatorId', o.id,
                           'empNo', o.emp_no,
                           'operatorName', o.name,
                           'section', o.section,
                           'operatorStatus', o.status
                   )
           ) AS operators
    FROM tbl_department td
             LEFT JOIN tbl_operator o ON o.department_id = td.id
    WHERE td.id = departmentId
    GROUP BY td.id, td.department_code, td.department_name, td.department_status;
END;