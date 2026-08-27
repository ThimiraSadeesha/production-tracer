DROP PROCEDURE IF EXISTS machine_history_get;
CREATE PROCEDURE machine_history_get(
    IN machine_id_val BIGINT
)
BEGIN

    SET time_zone = 'Asia/Colombo';

    WITH RECURSIVE
        quarter_slots AS (
            SELECT
                TIMESTAMP(DATE_FORMAT(
                        DATE_SUB(NOW(), INTERVAL (MOD(MINUTE(NOW()), 15) * 60 + SECOND(NOW())) SECOND),
                        '%Y-%m-%d %H:%i:00'
                          )) AS slot_start,
                TIMESTAMP(DATE_FORMAT(
                        DATE_ADD(
                                DATE_SUB(NOW(), INTERVAL (MOD(MINUTE(NOW()), 15) * 60 + SECOND(NOW())) SECOND),
                                INTERVAL 15 MINUTE
                        ),
                        '%Y-%m-%d %H:%i:00'
                          )) AS slot_end,
                0 AS n
            UNION ALL
            SELECT
                DATE_SUB(slot_start, INTERVAL 15 MINUTE),
                DATE_SUB(slot_end,   INTERVAL 15 MINUTE),
                n + 1
            FROM quarter_slots
            WHERE n < 95
        ),

        last_status_per_slot AS (
            SELECT
                qs.slot_start,
                -- Map raw log status → one of the 4 valid machine statuses
                CASE SUBSTRING_INDEX(
                        GROUP_CONCAT(ml.machine_status ORDER BY ml.started_at DESC),
                        ',', 1
                     )
                    WHEN 'Completed'   THEN 'Running'
                    WHEN 'MakeReady'  THEN 'Running'
                    WHEN 'Running'     THEN 'Running'
                    WHEN 'Paused'      THEN 'Paused'
                    WHEN 'Maintenance' THEN 'Maintenance'
                    ELSE 'Idle'
                    END AS machine_status
            FROM quarter_slots qs
                     JOIN tbl_machine_log ml
                          ON ml.machine_id = machine_id_val
                              AND ml.started_at <  qs.slot_end
                              AND (ml.ended_at IS NULL OR ml.ended_at >= qs.slot_start)
                              AND ml.started_at >= NOW() - INTERVAL 24 HOUR
            GROUP BY qs.slot_start
        )

    SELECT
        m.id           AS machine_id,
        m.machine_name,
        m.machine_code,
        m.status       AS current_status,
        COALESCE(
                JSON_ARRAYAGG(
                        JSON_OBJECT(
                                'slot',           qs.slot_start,
                                'machine_status', COALESCE(lsps.machine_status, 'Idle')
                        ) ORDER BY qs.slot_start ASC
                ),
                JSON_ARRAY()
        ) AS histories
    FROM tbl_machine m
             CROSS JOIN quarter_slots qs
             LEFT JOIN last_status_per_slot lsps ON lsps.slot_start = qs.slot_start
    WHERE m.id = machine_id_val
    GROUP BY m.id, m.machine_name, m.machine_code, m.status;

END;
