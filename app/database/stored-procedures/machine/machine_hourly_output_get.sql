DROP PROCEDURE IF EXISTS machine_hourly_output_get;
CREATE PROCEDURE machine_hourly_output_get(
    IN p_department_id BIGINT
)
BEGIN
    SELECT m.id            AS machine_id,
           m.machine_code  AS machine_code,
           m.machine_name  AS machine_name,
           m.status        AS machine_status,
           m.hourly_output AS hourly_output,
           mt.id           AS machine_type_id,
           mt.name         AS machine_type_name,
           u.id            AS unit_id,
           u.unit_code     AS unit_code,
           u.unit_name     AS unit_name,
           u.unit_symbol   AS unit_symbol,
           d.id            AS department_id,
           d.department_name AS department_name
    FROM tbl_machine m
             INNER JOIN tbl_machine_type mt ON m.machine_type_id = mt.id
             INNER JOIN tbl_process p ON mt.process_id = p.id
             INNER JOIN tbl_department d ON p.department_id = d.id
             LEFT JOIN tbl_unit u ON m.unit_of_machine_id = u.id
    WHERE m.deleted_at IS NULL
      AND m.status IN ('Idle', 'Running', 'Paused')
      AND (p_department_id IS NULL OR p_department_id = -1 OR d.id = p_department_id)
    ORDER BY m.machine_code ASC;
END;
