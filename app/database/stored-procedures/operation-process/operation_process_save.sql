DROP PROCEDURE IF EXISTS operation_process_save;
CREATE PROCEDURE operation_process_save(
    IN operation_process JSON,
    IN created_by_val VARCHAR(150)
)
BEGIN
    DECLARE total INT;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_err_code VARCHAR(20);
    DECLARE v_err_ord INT;
    DECLARE v_machine_ord INT;
    DECLARE v_order_exist_ord INT;
    DECLARE v_order_batch_ord INT;
    DECLARE v_items_exist_ord INT;
    DECLARE v_items_batch_ord INT;
    DECLARE v_order_dup_ord INT;
    DECLARE v_items_dup_ord INT;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            DROP TEMPORARY TABLE IF EXISTS tmp_ops_save;
            ROLLBACK;
        END;

    START TRANSACTION;

    SET total = JSON_LENGTH(operation_process);
    SET SESSION group_concat_max_len = GREATEST(@@session.group_concat_max_len, 1048576);

    DROP TEMPORARY TABLE IF EXISTS tmp_ops_save;
    CREATE TEMPORARY TABLE tmp_ops_save
    (
        ord              INT PRIMARY KEY,
        wo_item          BIGINT,
        wo_id            BIGINT,
        op_scope         VARCHAR(10) NOT NULL,
        process_id       BIGINT,
        machine_id       BIGINT,
        plan_date_time   VARCHAR(64),
        sequence_val     VARCHAR(64),
        status_val       VARCHAR(100),
        estimate_time    VARCHAR(64),
        planned_quantity VARCHAR(64),
        expected_value   VARCHAR(64),
        target_wo_id     BIGINT,
        INDEX idx_tmp_order (op_scope, wo_id, process_id, machine_id),
        INDEX idx_tmp_item (op_scope, wo_item, process_id, machine_id),
        INDEX idx_tmp_machine (machine_id),
        INDEX idx_tmp_target (target_wo_id)
    );

    INSERT INTO tmp_ops_save (ord, wo_item, wo_id, op_scope, process_id, machine_id,
                              plan_date_time, sequence_val, status_val, estimate_time,
                              planned_quantity, expected_value, target_wo_id)
    SELECT s.ord,
           s.wo_item,
           s.wo_id,
           s.op_scope,
           s.process_id,
           s.machine_id,
           s.plan_date_time,
           s.sequence_val,
           s.status_val,
           s.estimate_time,
           s.planned_quantity,
           s.expected_value,
           CASE
               WHEN s.op_scope = 'ORDER' AND s.wo_id IS NOT NULL AND s.wo_id > 0
                   THEN s.wo_id
               WHEN s.op_scope = 'ITEMS' AND s.wo_item IS NOT NULL AND s.wo_item > 0
                   THEN woi.work_order_id
               ELSE NULL
           END
    FROM (
        SELECT parsed.ord,
               CASE
                   WHEN parsed.wo_id_sql IS NOT NULL THEN NULL
                   ELSE parsed.wo_item_sql
               END                                              AS wo_item,
               parsed.wo_id_sql                                 AS wo_id,
               IF(parsed.wo_id_sql IS NOT NULL, 'ORDER', 'ITEMS') AS op_scope,
               parsed.process_id_sql                            AS process_id,
               parsed.machine_id_sql                            AS machine_id,
               parsed.plan_date_time,
               parsed.sequence_val,
               parsed.status_val,
               parsed.estimate_time,
               parsed.planned_quantity,
               parsed.expected_value
        FROM (
            SELECT jt.ord,
                   IF(wo_item_raw IS NULL OR wo_item_raw = '' OR LOWER(wo_item_raw) = 'null',
                      NULL, CAST(wo_item_raw AS SIGNED))        AS wo_item_sql,
                   IF(wo_id_raw IS NULL OR wo_id_raw = '' OR LOWER(wo_id_raw) = 'null',
                      NULL, CAST(wo_id_raw AS SIGNED))          AS wo_id_sql,
                   IF(process_id_raw IS NULL OR process_id_raw = '' OR LOWER(process_id_raw) = 'null',
                      NULL, CAST(process_id_raw AS SIGNED))     AS process_id_sql,
                   IF(machine_id_raw IS NULL OR machine_id_raw = '' OR LOWER(machine_id_raw) = 'null',
                      NULL, CAST(machine_id_raw AS SIGNED))     AS machine_id_sql,
                   NULLIF(jt.planDateTime, 'null')              AS plan_date_time,
                   NULLIF(jt.sequence_val, 'null')              AS sequence_val,
                   NULLIF(jt.status_val, 'null')                AS status_val,
                   NULLIF(jt.estimateTime, 'null')              AS estimate_time,
                   NULLIF(jt.plannedQuantity, 'null')           AS planned_quantity,
                   NULLIF(jt.expectedValue, 'null')             AS expected_value
            FROM (
                SELECT jt.ord,
                       TRIM(jt.workOrderItemId) AS wo_item_raw,
                       TRIM(jt.workOrderId)     AS wo_id_raw,
                       TRIM(jt.processId)       AS process_id_raw,
                       TRIM(jt.machineId)       AS machine_id_raw,
                       jt.planDateTime,
                       jt.sequence_val,
                       jt.status_val,
                       jt.estimateTime,
                       jt.plannedQuantity,
                       jt.expectedValue
                FROM JSON_TABLE(
                    operation_process,
                    '$[*]' COLUMNS (
                        ord FOR ORDINALITY,
                        workOrderItemId VARCHAR(64) PATH '$.workOrderItemId',
                        workOrderId VARCHAR(64) PATH '$.workOrderId',
                        processId VARCHAR(64) PATH '$.processId',
                        machineId VARCHAR(64) PATH '$.machineId',
                        planDateTime VARCHAR(64) PATH '$.planDateTime',
                        sequence_val VARCHAR(64) PATH '$.sequence',
                        status_val VARCHAR(100) PATH '$.status',
                        estimateTime VARCHAR(64) PATH '$.estimateTime',
                        plannedQuantity VARCHAR(64) PATH '$.plannedQuantity',
                        expectedValue VARCHAR(64) PATH '$.expectedValue'
                    )
                ) AS jt
            ) jt
        ) parsed
    ) s
    LEFT JOIN tbl_work_order_item woi ON woi.id = s.wo_item;

    SET v_machine_ord = (
        SELECT MIN(t.ord)
        FROM tmp_ops_save t
                 LEFT JOIN tbl_machine m ON m.id = t.machine_id
        WHERE t.machine_id IS NOT NULL
          AND t.machine_id > 0
          AND m.id IS NULL
    );

    SET v_order_exist_ord = (
        SELECT MIN(t.ord)
        FROM tmp_ops_save t
                 INNER JOIN tbl_operation_process op
                            ON op.operation_scope = 'ORDER'
                                AND op.work_order_id = t.wo_id
                                AND op.process_id = t.process_id
                                AND op.machine_id <=> t.machine_id
        WHERE t.op_scope = 'ORDER'
          AND t.wo_id IS NOT NULL AND t.wo_id > 0
          AND t.process_id IS NOT NULL AND t.process_id > 0
    );

    SET v_order_batch_ord = (
        SELECT MIN(d.ord)
        FROM (
            SELECT ord,
                   ROW_NUMBER() OVER (
                       PARTITION BY wo_id, process_id, machine_id
                       ORDER BY ord
                   ) AS rn
            FROM tmp_ops_save
            WHERE op_scope = 'ORDER'
              AND wo_id IS NOT NULL AND wo_id > 0
              AND process_id IS NOT NULL AND process_id > 0
        ) d
        WHERE d.rn > 1
    );

    SET v_items_exist_ord = (
        SELECT MIN(t.ord)
        FROM tmp_ops_save t
                 INNER JOIN tbl_operation_process op
                            ON op.operation_scope = 'ITEMS'
                                AND op.work_order_item_id = t.wo_item
                                AND op.process_id = t.process_id
                                AND op.machine_id <=> t.machine_id
        WHERE t.op_scope = 'ITEMS'
          AND t.wo_item IS NOT NULL AND t.wo_item > 0
          AND t.process_id IS NOT NULL AND t.process_id > 0
    );

    SET v_items_batch_ord = (
        SELECT MIN(d.ord)
        FROM (
            SELECT ord,
                   ROW_NUMBER() OVER (
                       PARTITION BY wo_item, process_id, machine_id
                       ORDER BY ord
                   ) AS rn
            FROM tmp_ops_save
            WHERE op_scope = 'ITEMS'
              AND wo_item IS NOT NULL AND wo_item > 0
              AND process_id IS NOT NULL AND process_id > 0
        ) d
        WHERE d.rn > 1
    );

    SET v_order_dup_ord = CASE
                              WHEN v_order_exist_ord IS NULL THEN v_order_batch_ord
                              WHEN v_order_batch_ord IS NULL THEN v_order_exist_ord
                              ELSE LEAST(v_order_exist_ord, v_order_batch_ord)
                          END;

    SET v_items_dup_ord = CASE
                              WHEN v_items_exist_ord IS NULL THEN v_items_batch_ord
                              WHEN v_items_batch_ord IS NULL THEN v_items_exist_ord
                              ELSE LEAST(v_items_exist_ord, v_items_batch_ord)
                          END;

    SET v_err_code = NULL;
    SET v_err_ord = NULL;

    IF v_machine_ord IS NOT NULL THEN
        SET v_err_ord = v_machine_ord;
        SET v_err_code = 'machine';
    END IF;

    IF v_order_dup_ord IS NOT NULL AND (v_err_ord IS NULL OR v_order_dup_ord < v_err_ord) THEN
        SET v_err_ord = v_order_dup_ord;
        SET v_err_code = 'order_dup';
    END IF;

    IF v_items_dup_ord IS NOT NULL AND (v_err_ord IS NULL OR v_items_dup_ord < v_err_ord) THEN
        SET v_err_ord = v_items_dup_ord;
        SET v_err_code = 'items_dup';
    END IF;

    IF v_err_code = 'machine' THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'The provided machineId does not exist.';
    ELSEIF v_err_code = 'order_dup' THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT =
                    'This work order is already scheduled for this process (and machine) at ORDER scope.';
    ELSEIF v_err_code = 'items_dup' THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT =
                    'This work order item is already scheduled for this process (and machine) at ITEMS scope.';
    END IF;

    INSERT INTO tbl_operation_process (plan_date_time,
                                       sequence,
                                       status,
                                       estimate_time,
                                       machine_id,
                                       work_order_item_id,
                                       work_order_id,
                                       operation_scope,
                                       process_id,
                                       planned_quantity,
                                       expected_value,
                                       created_by)
    SELECT t.plan_date_time,
           t.sequence_val,
           t.status_val,
           t.estimate_time,
           t.machine_id,
           t.wo_item,
           t.wo_id,
           t.op_scope,
           t.process_id,
           t.planned_quantity,
           t.expected_value,
           created_by_val
    FROM tmp_ops_save t
    ORDER BY t.ord;

    SET inserted_ids = (
        SELECT CAST(CONCAT('[', IFNULL(GROUP_CONCAT(x.id ORDER BY x.ord SEPARATOR ','), ''), ']') AS JSON)
        FROM (
            SELECT t.ord, MAX(op.id) AS id
            FROM tmp_ops_save t
                     INNER JOIN tbl_operation_process op
                                ON op.operation_scope = t.op_scope
                                    AND op.process_id <=> t.process_id
                                    AND op.machine_id <=> t.machine_id
                                    AND (
                                        (t.op_scope = 'ORDER' AND op.work_order_id <=> t.wo_id)
                                            OR (t.op_scope = 'ITEMS' AND op.work_order_item_id <=> t.wo_item)
                                        )
            GROUP BY t.ord
        ) x
    );

    UPDATE tbl_work_order_item woi
        INNER JOIN tmp_ops_save t ON t.op_scope = 'ORDER' AND t.wo_id = woi.work_order_id
    SET woi.item_status = 'Scheduled';

    UPDATE tbl_work_order_item woi
        INNER JOIN tmp_ops_save t ON t.op_scope = 'ITEMS' AND t.wo_item = woi.id
    SET woi.item_status = 'Scheduled';

    UPDATE tbl_work_order wo
        INNER JOIN tmp_ops_save t ON t.target_wo_id = wo.id
    SET wo.status = 'Scheduled';

    DROP TEMPORARY TABLE IF EXISTS tmp_ops_save;

    COMMIT;
    SELECT total        AS SavedTotalOperationProcess,
           inserted_ids AS InsertedIds;

END;
