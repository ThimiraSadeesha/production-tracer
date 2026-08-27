DROP PROCEDURE IF EXISTS get_machine_time_line;
CREATE PROCEDURE get_machine_time_line(
    IN p_machine_id BIGINT,
    IN p_from_date DATE,
    IN p_to_date DATE
)
BEGIN

    DECLARE v_from_date DATE;
    DECLARE v_to_date DATE;

    SET v_from_date = IFNULL(p_from_date, DATE_SUB(CURDATE(), INTERVAL 2 DAY));
    SET v_to_date = IFNULL(p_to_date, DATE_ADD(CURDATE(), INTERVAL 3 DAY));

    SELECT m.id                                     AS machine_id,
           m.machine_code,
           m.machine_name,
           m.status,
           (SELECT JSON_OBJECTAGG(plan_date, day_operations)
            FROM (SELECT DATE(op.plan_date_time) AS plan_date,
                         JSON_ARRAYAGG(
                                 JSON_OBJECT(
                                         'operation_id', op.id,
                                         'plan_date_time', op.plan_date_time,
                                         'estimate_time', op.estimate_time,
                                         'started_time', op.started_date_time,
                                         'completed_time', op.completed_date_time,
                                         'status', op.status,
                                         'pause_reason',
                                         CASE
                                             WHEN LOWER(op.status) = 'paused' THEN (SELECT thr.reason
                                                                                    FROM tbl_operation_log ol
                                                                                             JOIN tbl_hold_reason thr ON ol.reason = thr.id
                                                                                    WHERE ol.job_operation_id = op.id
                                                                                      AND ol.resume_date_time IS NULL
                                                                                      AND ol.reason IS NOT NULL
                                                                                    ORDER BY ol.pause_date_time DESC
                                                                                    LIMIT 1)
                                             ELSE NULL
                                             END,
                                         'inventory_item', woi.inventory_item,
                                         'pp_number', woi.pp_number,
                                         'size', woi.size,
                                         'item_quantity', COALESCE(woi.quantity, wo.quantity),
                                         'work_order_id', wo.id,
                                         'work_order_reference', wo.work_order_reference,
                                         'work_order_type', wo.work_order_type,
                                         'no_of_breakdowns', wo.no_of_breakdowns,
                                         'customer_name', wo.customer_name,
                                         'product_design', wo.product_design,
                                         'department_id', wo.department_id
                                 )
                         )                       AS day_operations
                  FROM tbl_operation_process op
                           LEFT JOIN tbl_work_order_item woi ON woi.id = op.work_order_item_id
                           LEFT JOIN tbl_work_order wo ON wo.id = IFNULL(woi.work_order_id, op.work_order_id)
                  WHERE op.machine_id = p_machine_id
                    AND op.plan_date_time >= v_from_date
                    AND op.plan_date_time <= DATE_ADD(v_to_date, INTERVAL 1 DAY)
                  GROUP BY DATE(op.plan_date_time)
                  ORDER BY plan_date) AS dated_ops) AS operations
    FROM tbl_machine m
    WHERE (p_machine_id IS NULL OR m.id = p_machine_id)
    ORDER BY m.id;

END;
