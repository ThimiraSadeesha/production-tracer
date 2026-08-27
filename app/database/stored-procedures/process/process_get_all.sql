DROP PROCEDURE IF EXISTS process_getAll;
CREATE PROCEDURE process_getAll()
BEGIN
    SELECT tp.id           AS processId,
           tp.code         AS processCode,
           tp.name         AS processName,
           tp.status       AS processStatus,
           tp.is_required  AS isRequired,
           tp.department_id AS departmentId
    FROM tbl_process tp;
END;