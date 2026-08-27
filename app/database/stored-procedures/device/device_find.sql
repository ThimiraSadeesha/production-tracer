DROP PROCEDURE IF EXISTS device_find;
CREATE PROCEDURE device_find(
    IN android_id_val VARCHAR(150),
    IN process_id_val BIGINT,
    IN status_val VARCHAR(50),
#     IN device_type_val VARCHAR(50),
    IN machine_id_val BIGINT,
    IN p_cursor BIGINT,
    IN p_limit INT
)
BEGIN
    SET collation_connection = 'utf8mb4_unicode_ci';


    IF p_limit IS NULL OR p_limit <= 0 THEN
        SET p_limit = 20;
    END IF;

    SELECT d.id                 AS deviceId,
           d.app_version        AS appVersion,
           d.last_seen_at       AS lastSeenAt,
           d.status             AS deviceStatus,
           d.device_type        AS deviceType,
           d.android_id         AS androidId,
           d.machine_id         AS machineId,
           d.process_id         AS processId,
           m.machine_code       AS machineCode,
           m.machine_name       AS machineName,
           dept.department_code AS departmentCode,
           dept.id              AS departmentId,
           p.name               AS processName
    FROM tbl_device d
             LEFT JOIN tbl_machine m ON d.machine_id = m.id
             LEFT JOIN tbl_process p ON d.process_id = p.id
             LEFT JOIN tbl_department dept ON p.department_id = dept.id
    WHERE (android_id_val IS NULL OR android_id_val = '' OR d.android_id LIKE CONCAT('%', android_id_val, '%'))
      AND (status_val IS NULL OR status_val = '' OR d.status LIKE CONCAT('%', status_val, '%'))
#       AND (device_type_val IS NULL OR device_type_val = '' OR d.device_type = device_type_val)
      AND (process_id_val IS NULL OR process_id_val = -1 OR d.process_id = process_id_val)
      AND (machine_id_val IS NULL OR machine_id_val = -1 OR d.machine_id = machine_id_val)
      AND (p_cursor = -1 OR d.id < p_cursor)
    ORDER BY d.id DESC
    LIMIT p_limit;

END;
