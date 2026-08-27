DROP PROCEDURE IF EXISTS device_identify;
CREATE PROCEDURE device_identify(
    IN android_id_val VARCHAR(150)
)
BEGIN
    -- heartbeat
    UPDATE tbl_device
    SET last_seen_at = NOW()
    WHERE android_id = android_id_val;

    SELECT d.id                                                 AS deviceId,
           d.android_id                                         AS androidId,
           d.device_type                                        AS deviceType,
           d.status                                             AS deviceStatus,
           d.app_version                                        AS appVersion,
           d.last_seen_at                                       AS lastSeenAt,
           d.machine_id                                         AS machineId,
           d.process_id                                         AS processId,
           p.name                                               AS processName,
           p.code                                               AS processCode,
           p.status                                             AS processStatus,
           dept.id                                              AS departmentId,
           dept.department_name                                 AS departmentName,
           dept.department_code                                 AS departmentCode,
           (SELECT JSON_ARRAYAGG(
                           JSON_OBJECT(
                                   'id', machine_rows.id,
                                   'machineCode', machine_rows.machine_code,
                                   'machineName', machine_rows.machine_name,
                                   'status', machine_rows.status
                           )
                   )
            FROM (SELECT DISTINCT m2.id,
                                  m2.machine_code,
                                  m2.machine_name,
                                  m2.status
                  FROM tbl_device d2
                           INNER JOIN tbl_machine m2 ON m2.id = d2.machine_id
                  WHERE d2.android_id = android_id_val
                    AND d2.deleted_at IS NULL
                    AND m2.deleted_at IS NULL) AS machine_rows) AS machines
    FROM tbl_device d
             LEFT JOIN tbl_process p
                       ON p.id = d.process_id
             LEFT JOIN tbl_department dept
                       ON dept.id = p.department_id
    WHERE d.android_id = android_id_val
    ORDER BY d.id DESC
    LIMIT 1;

END;
