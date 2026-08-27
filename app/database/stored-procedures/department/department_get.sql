DROP PROCEDURE IF EXISTS department_get;
CREATE PROCEDURE department_get(
    IN departmentId INT
)
BEGIN
    SELECT d.id                AS id,
           d.department_code   AS departmentCode,
           d.department_name   AS departmentName,
           d.department_status AS status,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'processId', ds.id,
                           'processName', ds.name,
                           'processStatus', ds.status,
                           'isRequired', ds.is_required
                   )
           )                   AS processes
    FROM tbl_department d
             LEFT JOIN tbl_process ds ON d.id = ds.department_id
    WHERE d.id = departmentId
    GROUP BY d.id, d.department_code, d.department_name;

END
