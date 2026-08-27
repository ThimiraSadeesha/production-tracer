DROP PROCEDURE IF EXISTS operation_process_get;
CREATE PROCEDURE operation_process_get(
    IN p_id BIGINT
)
BEGIN
    SELECT op.id                  AS id,
           op.sequence            AS sequence,
           op.plan_date_time      AS planDateTime,
           op.started_date_time   AS startedDateTime,
           op.completed_date_time AS completedDateTime,
           op.planned_quantity    AS plannedQuantity,
           op.expected_value      AS expectedValue,
           op.actual_value        AS actualValue,
           op.status              AS status,
           op.estimate_time       AS estimateTime,
           op.operation_scope     AS operationScope,
           op.is_redo             AS isRedo,
           op.remark              AS remark,
           op.machine_id          AS machineId,
           m.machine_code         AS machineCode,
           m.machine_name         AS machineName,
           op.process_id          AS processId,
           p.code                 AS processCode,
           p.name                 AS processName,
           op.work_order_id       AS workOrderId,
           wo.work_order_reference AS workOrderReference,
           op.work_order_item_id  AS workOrderItemId,
           woi.item_reference     AS itemReference,
           woi.item_name          AS itemName,
           IFNULL(
                   (SELECT JSON_ARRAYAGG(
                                   JSON_OBJECT(
                                           'id', opa.id,
                                           'attributeId', opa.attribute_id,
                                           'attributeName', pa.attribute_name,
                                           'value', opa.value
                                   )
                           )
                    FROM tbl_operation_process_attributes opa
                             LEFT JOIN tbl_process_attributes pa ON pa.id = opa.attribute_id
                    WHERE opa.operation_process_id = op.id),
                   JSON_ARRAY()
           )                      AS attributes
    FROM tbl_operation_process op
             LEFT JOIN tbl_machine m ON m.id = op.machine_id
             LEFT JOIN tbl_process p ON p.id = op.process_id
             LEFT JOIN tbl_work_order wo ON wo.id = op.work_order_id
             LEFT JOIN tbl_work_order_item woi ON woi.id = op.work_order_item_id
    WHERE op.id = p_id
      AND op.deleted_at IS NULL;
END;
