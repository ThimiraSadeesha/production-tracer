DROP PROCEDURE IF EXISTS operation_process_find;
CREATE PROCEDURE operation_process_find(
    IN p_machine_id BIGINT,
    IN p_process_id BIGINT,
    IN p_planned_date VARCHAR(255),
    IN p_work_order_id VARCHAR(255),
    IN p_status VARCHAR(100),
    IN p_cursor BIGINT,
    IN p_limit INT
)
BEGIN
    IF p_limit IS NULL OR p_limit <= 0 THEN
        SET p_limit = 20;
    END IF;

    SELECT op.id                   AS id,
           op.sequence             AS sequence,
           op.plan_date_time       AS planDateTime,
           op.status               AS status,
           op.operation_scope      AS operationScope,
           op.planned_value        AS plannedValue,
           op.machine_id           AS machineId,
           m.machine_name          AS machineName,
           op.process_id           AS processId,
           p.name                  AS processName,
           op.work_order_id        AS workOrderId,
           wo.work_order_reference AS workOrderReference,
           op.work_order_item_id   AS workOrderItemId,
           woi.item_reference      AS itemReference
    FROM tbl_operation_process op
             LEFT JOIN tbl_machine m ON m.id = op.machine_id
             LEFT JOIN tbl_process p ON p.id = op.process_id
             LEFT JOIN tbl_work_order wo ON wo.id = op.work_order_id
             LEFT JOIN tbl_work_order_item woi ON woi.id = op.work_order_item_id
    WHERE op.deleted_at IS NULL
      AND (p_machine_id IS NULL OR p_machine_id <= 0 OR op.machine_id = p_machine_id)
      AND (p_process_id IS NULL OR p_process_id <= 0 OR op.process_id = p_process_id)
      AND (p_planned_date IS NULL OR p_planned_date = '' OR DATE(op.plan_date_time) = p_planned_date)
      AND (p_work_order_id IS NULL OR p_work_order_id = '' OR op.work_order_id = p_work_order_id)
      AND (p_status IS NULL OR p_status = '' OR op.status = p_status)
      AND (p_cursor = -1 OR op.id < p_cursor)
    ORDER BY op.id DESC
    LIMIT p_limit;
END;
