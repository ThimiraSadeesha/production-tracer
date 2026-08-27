DROP PROCEDURE IF EXISTS shift_getAll;
CREATE PROCEDURE shift_getAll()
BEGIN
    SELECT ts.id              AS shiftId,
           ts.start_date      AS startDate,
           ts.shift           AS shift,
           ts.end_date        AS endDate,
           ts.start_time      AS startTime,
           ts.end_time        AS endTime,
           ts.department_id   AS departmentId,
           td.department_name AS departmentName,
           td.department_code AS departmentCode,
           ts.status          AS shiftStatus
    FROM tbl_shift ts
             INNER JOIN tbl_department td ON td.id = ts.department_id;

END;
