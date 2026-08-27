DROP PROCEDURE IF EXISTS machine_type_get;
CREATE PROCEDURE machine_type_get(IN machine_type_id INT)
BEGIN
    SELECT mt.id                  AS machineTypeId,
           mt.code,
           mt.name,
           mt.machine_type_status As status,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'id', m.id,
                           'machineCode', m.machine_code,
                           'machineName', m.machine_name,
                           'description', m.description,
                           'capabilities', m.capabilities,
                           'status', m.status,
                           'hourlyOutput', m.hourly_output,
                           'makeReadyTime', m.make_ready_time,
                           'unitId', m.unit_of_machine_id,
                           'unit', u.unit_code,
                           'unitName', u.unit_name
                   )
           )                      AS machines
    FROM tbl_machine_type mt
             LEFT JOIN tbl_machine m ON mt.id = m.machine_type_id
             LEFT JOIN tbl_unit u ON m.unit_of_machine_id = u.id
    WHERE mt.id = machine_type_id
    GROUP BY mt.id;
END;

