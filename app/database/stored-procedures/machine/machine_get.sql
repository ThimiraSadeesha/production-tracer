DROP PROCEDURE IF EXISTS machine_get;
CREATE PROCEDURE machine_get(
    IN machine_id_val BIGINT
)
BEGIN
    SELECT m.id                                                AS machine_id,
           m.machine_code                                      AS machine_code,
           m.machine_name                                      AS machine_name,
           m.description                                       AS description,

           m.status                                            AS machine_status,
           m.hourly_output                                     AS hourly_output,
           m.make_ready_time                                   AS make_ready_time,
           m.machine_type_id                                   AS machine_type_id,
           mt.name                                             AS machine_type_name,
           JSON_OBJECT(
                   'unit_id', u.id,
                   'unit_code', u.unit_code,
                   'unit_name', u.unit_name,
                   'unit_symbol', u.unit_symbol
           )                                                   AS unit,
           IFNULL((SELECT JSON_ARRAYAGG(
                                  JSON_OBJECT(
                                          'reason_id', mhr.id,
                                          'reason', hr.reason
                                  )
                          )
                   FROM tbl_machine_hold_reason mhr
                            JOIN tbl_hold_reason hr ON mhr.hold_reason_id = hr.id
                   WHERE mhr.machine_id = m.id AND  hr.id <> 13) , JSON_ARRAY()) AS hold_reasons,

           (SELECT JSON_OBJECT(
                           'total_working_hours',
                           ROUND(IFNULL(SUM(
                                                CASE
                                                    WHEN ml.machine_status IN ('Completed', 'Running', 'MakeReady')
                                                        THEN ABS(TIMESTAMPDIFF(SECOND, ml.started_at, IFNULL(ml.ended_at, NOW())))
                                                    ELSE 0
                                                    END
                                        ) / 3600, 0), 2),

                           'total_up_time_hours',
                           ROUND(IFNULL(SUM(
                                                CASE
                                                    WHEN ml.machine_status NOT IN ('Maintenance', 'Paused')
                                                        THEN ABS(TIMESTAMPDIFF(SECOND, ml.started_at, IFNULL(ml.ended_at, NOW())))
                                                    ELSE 0
                                                    END
                                        ) / 3600, 0), 2),

                           'total_down_time_hours',
                           ROUND(IFNULL(SUM(
                                                CASE
                                                    WHEN ml.machine_status IN ('Paused', 'Maintenance')
                                                        THEN ABS(TIMESTAMPDIFF(SECOND, ml.started_at, IFNULL(ml.ended_at, NOW())))
                                                    ELSE 0
                                                    END
                                        ) / 3600, 0), 2),

                           'efficiency_percentage',
                           ROUND(
                                   (
                                       SUM(
                                               CASE
                                                   WHEN ml.machine_status = 'Running'
                                                       THEN ABS(TIMESTAMPDIFF(SECOND, ml.started_at, IFNULL(ml.ended_at, NOW())))
                                                   END
                                       )
                                           /
                                       NULLIF(
                                               SUM(ABS(TIMESTAMPDIFF(SECOND, ml.started_at, IFNULL(ml.ended_at, NOW())))),
                                               0
                                       )
                                       ) * 100,
                                   2),

                           'expected_output',
                           ROUND(IFNULL((
                                            m.hourly_output *
                                            (
                                                SUM(
                                                        CASE
                                                            WHEN ml.machine_status = 'Running'
                                                                THEN ABS(TIMESTAMPDIFF(SECOND, ml.started_at, IFNULL(ml.ended_at, NOW())))
                                                            END
                                                ) / 3600
                                                )
                                            ), 0), 2)
                   )
            FROM tbl_machine_log ml
            WHERE ml.machine_id = m.id)                        AS performance,
           IFNULL(
                   (SELECT JSON_OBJECT(
                                   'operation_process_id', op.id,
                                   'operation_status',     op.status,
                                   'operation_scope',      op.operation_scope,
                                   'sequence',             op.sequence,
                                   'plan_date_time',       op.plan_date_time,
                                   'started_date_time',    op.started_date_time,
                                   'planned_quantity',     op.planned_quantity,
                                   'estimate_time',        op.estimate_time,
                                   'machine_log', (
                                       SELECT JSON_OBJECT(
                                                      'machine_log_id',   ml.id,
                                                      'machine_status',   ml.machine_status,
                                                      'started_at',       ml.started_at,
                                                      'ended_at',         ml.ended_at,
                                                      'reason',           ml.reason
                                              )
                                       FROM tbl_machine_log ml
                                       WHERE ml.machine_id       = m.id
                                         AND ml.job_operation_id = op.id
                                         AND ml.ended_at         IS NULL
                                       ORDER BY ml.started_at DESC
                                       LIMIT 1
                                   ),
                                   'make_ready', (
                                       SELECT JSON_ARRAYAGG(
                                                      JSON_OBJECT(
                                                              'make_ready_id',  mr.id,
                                                              'machine_id',     mr.machine_id,
                                                              'operation_process_id',     mr.operation_process_id,
                                                              'status',         mr.status,
                                                              'started_at',     mr.started_at,
                                                              'ended_at',       mr.ended_at,
                                                              'attributes', (
                                                                  SELECT JSON_ARRAYAGG(
                                                                                 JSON_OBJECT(
                                                                                         'attribute_id',   mra.attribute_id,
                                                                                         'value',          mra.value,
                                                                                         'attribute_name', pa.attribute_name
                                                                                 )
                                                                         )
                                                                  FROM tbl_make_ready_attributes mra
                                                                          left JOIN tbl_process_attributes pa ON pa.id = mra.attribute_id
                                                                  WHERE mra.make_ready_id = mr.id
                                                              )

                                                      )
                                              )
                                       FROM tbl_make_ready mr
                                       WHERE mr.machine_id = machine_id_val AND mr.status IN ('Processing', 'Paused')
                                   ),
                                   'operation_attributes', (
                                       SELECT JSON_ARRAYAGG(
                                                      JSON_OBJECT(
                                                              'attribute_id',   opa.attribute_id,
                                                              'value',          opa.value,
                                                              'source',         opa.source
                                                      )
                                              )
                                       FROM tbl_operation_process_attributes opa
                                                JOIN tbl_process_attributes pa
                                                     ON pa.id = opa.attribute_id
                                       WHERE opa.operation_process_id = op.id
                                         AND opa.source               = 'Operation'
                                   ),
                                   'work_order', JSON_OBJECT(
                                           'work_order_id',        COALESCE(wo.id,                   wo_item.id),
                                           'work_order_reference', COALESCE(wo.work_order_reference, wo_item.work_order_reference),
                                           'product_design',       COALESCE(wo.product_design,       wo_item.product_design),
                                           'product_type',         COALESCE(wo.product_type,         wo_item.product_type),
                                           'product_sample',       COALESCE(wo.product_sample,       wo_item.product_sample),
                                           'quantity',             COALESCE(wo.quantity,             wo_item.quantity),
                                           'status',               COALESCE(wo.status,               wo_item.status),
                                           'delivery_date',        COALESCE(wo.delivery_date,        wo_item.delivery_date),
                                           'customer_name',        COALESCE(wo.customer_name,        wo_item.customer_name)
                                                 ),
                                   'work_order_item', CASE
                                                          WHEN op.operation_scope = 'ITEMS' AND woi.id IS NOT NULL
                                                              THEN JSON_OBJECT(
                                                                  'work_order_item_id',   woi.id,
                                                                  'work_order_id',        woi.work_order_id,
                                                                  'pp_number',            woi.pp_number,
                                                                  'inventory_item',       woi.inventory_item,
                                                                  'size',                 woi.size,
                                                                  'quantity',             woi.quantity,
                                                                  'produced_quantity',    woi.produced_quantity,
                                                                  'cut_quantity',         woi.cut_quantity,
                                                                  'pack_quantity',        woi.pack_quantity,
                                                                  'item_status',          woi.item_status,
                                                                  'remarks',              woi.remarks
                                                                   )
                                                          ELSE NULL
                                       END
                           )
                    FROM tbl_operation_process op
                             LEFT JOIN tbl_work_order wo
                                       ON  wo.id              = op.work_order_id
                                           AND op.operation_scope = 'ORDER'
                                           AND wo.deleted_at      IS NULL
                             LEFT JOIN tbl_work_order_item woi
                                       ON  woi.id             = op.work_order_item_id
                                           AND op.operation_scope = 'ITEMS'
                                           AND woi.deleted_at     IS NULL
                             LEFT JOIN tbl_work_order wo_item
                                       ON  wo_item.id         = woi.work_order_id
                                           AND op.operation_scope = 'ITEMS'
                                           AND wo_item.deleted_at IS NULL
                    WHERE op.machine_id = m.id
                      AND op.status     IN('Processing','MakeReady','Paused','MakeReadyPaused')
                      AND (
                        (op.operation_scope = 'ORDER' AND wo.id  IS NOT NULL)
                            OR (op.operation_scope = 'ITEMS' AND woi.id IS NOT NULL)
                        )
                    LIMIT 1),
                   JSON_OBJECT()) AS current_operation

    FROM tbl_machine m
             LEFT JOIN tbl_machine_type mt ON m.machine_type_id = mt.id
             LEFT JOIN tbl_unit u ON m.unit_of_machine_id = u.id
    WHERE m.id = machine_id_val;

END;
