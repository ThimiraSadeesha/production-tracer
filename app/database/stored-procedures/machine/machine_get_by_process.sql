DROP PROCEDURE IF EXISTS machine_get_by_process;
CREATE PROCEDURE machine_get_by_process(
    IN process_id BIGINT
)
BEGIN

    SELECT p.id   AS processId,
           p.name AS processName,
           p.code AS processCode,
           JSON_ARRAYAGG(
                   JSON_OBJECT(
                           'id', mt.id,
                           'code', mt.code,
                           'name', mt.name,
                           'machines',
                           COALESCE(
                                   (SELECT JSON_ARRAYAGG(
                                                   JSON_OBJECT(
                                                           'id', m.id,
                                                           'machineCode', m.machine_code,
                                                           'machineName', m.machine_name,
                                                           'status', m.status,
                                                           'hourlyOutput', m.hourly_output
                                                   )
                                           )
                                    FROM tbl_machine m
                                    WHERE m.machine_type_id = mt.id),
                                   JSON_ARRAY()
                           )
                   )
           )      AS machine_types
    FROM tbl_process p
             JOIN tbl_machine_type mt ON mt.process_id = p.id
    WHERE p.id = process_id
    GROUP BY p.id, p.name, p.code;

END;