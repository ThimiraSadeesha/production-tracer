DROP PROCEDURE IF EXISTS shift_get_by_department;
CREATE PROCEDURE shift_get_by_department(
    IN departmentId INT
)
BEGIN
    SELECT td.id                AS departmentId,
           td.department_code   AS departmentCode,
           td.department_name   AS departmentName,
           td.department_status AS departmentStatus,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'shiftId', ts.id,
                           'shiftName', ts.shift,
                           'startDate', ts.start_date,
                           'endDate', ts.end_date,
                           'startTime', ts.start_time,
                           'endTime', ts.end_time,
                           'shiftStatus', ts.status
                   )
           )                    AS shifts
    FROM tbl_department td
             LEFT JOIN tbl_shift ts ON ts.department_id = td.id
    WHERE td.id = departmentId
    GROUP BY td.id, td.department_code, td.department_name, td.department_status;
END;
