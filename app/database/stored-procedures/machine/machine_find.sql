DROP PROCEDURE IF EXISTS machine_find;
CREATE PROCEDURE machine_find(
    IN machine_code_val VARCHAR(100),
    IN machine_name_val VARCHAR(255),
    IN machine_type_id_val BIGINT,
    IN department_id_val BIGINT,
    IN process_id_val BIGINT, -- NEW PARAM
    IN p_cursor BIGINT,
    IN p_limit INT
)
BEGIN
    -- Default limit
#     IF p_limit IS NULL OR p_limit <= 0 THEN
#         SET p_limit = 20;
#     END IF;

    SELECT m.id              AS machineId,
           m.machine_code    AS machineCode,
           m.machine_name    AS machineName,
           m.description     AS description,

           m.status          AS machineStatus,
           m.hourly_output   AS hourlyOutput,
           m.make_ready_time AS makeReadyTime,
           mt.name           AS machineTypeName,
           p.id              AS processId,   -- OPTIONAL (useful)
           p.name            AS processName, -- OPTIONAL
           u.unit_name       AS unitName,
           u.id              AS unitId,
           d.id              AS departmentId,
           d.department_name AS departmentName
    FROM tbl_machine m
             INNER JOIN tbl_machine_type mt ON m.machine_type_id = mt.id
             INNER JOIN tbl_process p ON mt.process_id = p.id
             INNER JOIN tbl_department d ON p.department_id = d.id
             LEFT JOIN tbl_unit u ON m.unit_of_machine_id = u.id
    WHERE (machine_code_val IS NULL OR machine_code_val = '' OR machine_code_val = 'null'
        OR m.machine_code LIKE CONCAT('%', machine_code_val, '%'))
      AND (machine_name_val IS NULL OR machine_name_val = '' OR machine_name_val = 'null'
        OR m.machine_name LIKE CONCAT('%', machine_name_val, '%'))
      AND (machine_type_id_val IS NULL OR machine_type_id_val = -1 OR m.machine_type_id = machine_type_id_val)
      AND (department_id_val IS NULL OR department_id_val = -1 OR d.id = department_id_val)
      AND (process_id_val IS NULL OR process_id_val = -1 OR p.id = process_id_val)
      AND (p_cursor IS NULL OR p_cursor = -1 OR m.id < p_cursor)
    ORDER BY m.id DESC;
#     LIMIT p_limit;
END;
