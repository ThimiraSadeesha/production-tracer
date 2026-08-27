DROP PROCEDURE IF EXISTS device_get_by_id;
CREATE PROCEDURE device_get_by_id(
    IN device_id_val BIGINT
)
BEGIN
    SELECT d.id           AS deviceId,
           d.app_version  AS appVersion,
           d.last_seen_at AS lastSeenAt,
           d.status       AS deviceStatus,
           d.device_type  AS deviceType,
           d.android_id   AS androidId,
           d.machine_id   AS machineId,
           d.process_id   AS processId,
           (
               SELECT JSON_ARRAYAGG(
                              JSON_OBJECT(
                                      'id', machine_rows.id,
                                      'machineCode', machine_rows.machine_code,
                                      'machineName', machine_rows.machine_name,
                                      'description', machine_rows.description,
                                      'capabilities', machine_rows.capabilities,
                                      'status', machine_rows.status,
                                      'hourlyOutput', machine_rows.hourly_output,
                                      'makeReadyTime', machine_rows.make_ready_time
                              )
                      )
               FROM (SELECT DISTINCT m2.id,
                                     m2.machine_code,
                                     m2.machine_name,
                                     m2.description,
                                     m2.capabilities,
                                     m2.status,
                                     m2.hourly_output,
                                     m2.make_ready_time
                     FROM tbl_device d2
                              INNER JOIN tbl_machine m2 ON d2.machine_id = m2.id
                     WHERE d2.android_id = d.android_id
                       AND d2.deleted_at IS NULL
                       AND m2.deleted_at IS NULL) AS machine_rows
           )              AS machines,
           (
               SELECT JSON_ARRAYAGG(
                              JSON_OBJECT(
                                      'id', process_rows.id,
                                      'name', process_rows.name,
                                      'code', process_rows.code,
                                      'status', process_rows.status
                              )
                      )
               FROM (SELECT DISTINCT p2.id,
                                     p2.name,
                                     p2.code,
                                     p2.status
                     FROM tbl_device d3
                              INNER JOIN tbl_process p2 ON d3.process_id = p2.id
                     WHERE d3.android_id = d.android_id
                       AND d3.deleted_at IS NULL
                       AND p2.deleted_at IS NULL) AS process_rows
           )              AS processes
    FROM tbl_device d
    WHERE d.id = device_id_val
      AND d.deleted_at IS NULL;
END;