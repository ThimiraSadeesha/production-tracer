DROP PROCEDURE IF EXISTS operation_process_save;
CREATE PROCEDURE operation_process_save(
    IN operation_process JSON,
    IN created_by_val VARCHAR(150)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_elem JSON;
    DECLARE v_wo_item BIGINT;
    DECLARE v_wo_id BIGINT;
    DECLARE v_scope VARCHAR(10);
    DECLARE v_process_id BIGINT;
    DECLARE v_machine_id BIGINT;
    DECLARE v_scope_raw VARCHAR(20);
    DECLARE v_inserted_id BIGINT;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    DROP TEMPORARY TABLE IF EXISTS tmp_sched_wo;
    DROP TEMPORARY TABLE IF EXISTS tmp_sched_woi;
    CREATE TEMPORARY TABLE tmp_sched_wo (id BIGINT PRIMARY KEY) ENGINE = MEMORY;
    CREATE TEMPORARY TABLE tmp_sched_woi (id BIGINT PRIMARY KEY) ENGINE = MEMORY;

    SET total = JSON_LENGTH(operation_process);

    START TRANSACTION;

    WHILE i < total
        DO

            SET v_elem = JSON_EXTRACT(operation_process, CONCAT('$[', i, ']'));

            SET v_wo_item = CAST(NULLIF(JSON_VALUE(v_elem, '$.workOrderItemId'), 0) AS SIGNED);
            SET v_wo_id = CAST(NULLIF(JSON_VALUE(v_elem, '$.workOrderId'), 0) AS SIGNED);
            SET v_process_id = CAST(NULLIF(JSON_VALUE(v_elem, '$.processId'), 0) AS SIGNED);
            SET v_machine_id = CAST(NULLIF(JSON_VALUE(v_elem, '$.machineId'), 0) AS SIGNED);

            SET v_scope_raw = UPPER(TRIM(JSON_VALUE(v_elem, '$.operationScope')));
            SET v_scope = CASE
                              WHEN v_scope_raw IN ('ORDER', 'ITEMS') THEN v_scope_raw
                              WHEN v_wo_id IS NOT NULL AND v_wo_id > 0 THEN 'ORDER'
                              ELSE 'ITEMS'
                END;

            IF v_scope = 'ORDER' THEN
                SET v_wo_item = NULL;
            ELSE
                SET v_wo_id = NULL;
            END IF;

            IF v_machine_id IS NOT NULL
                AND NOT EXISTS (SELECT 1 FROM tbl_machine WHERE id = v_machine_id) THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'The provided machineId does not exist.';
            END IF;

            IF v_scope = 'ORDER'
                AND v_wo_id IS NOT NULL
                AND v_process_id IS NOT NULL
                AND EXISTS (SELECT 1
                            FROM tbl_operation_process
                            WHERE work_order_id = v_wo_id
                              AND process_id = v_process_id
                              AND operation_scope = 'ORDER'
                              AND machine_id <=> v_machine_id) THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT =
                            'This work order is already scheduled for this process (and machine) at ORDER scope.';
            END IF;

            IF v_scope = 'ITEMS'
                AND v_wo_item IS NOT NULL
                AND v_process_id IS NOT NULL
                AND EXISTS (SELECT 1
                            FROM tbl_operation_process
                            WHERE work_order_item_id = v_wo_item
                              AND process_id = v_process_id
                              AND operation_scope = 'ITEMS'
                              AND machine_id <=> v_machine_id) THEN
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
            VALUES (JSON_VALUE(v_elem, '$.planDateTime'),
                    JSON_VALUE(v_elem, '$.sequence'),
                    JSON_VALUE(v_elem, '$.status'),
                    JSON_VALUE(v_elem, '$.estimateTime'),
                    v_machine_id,
                    v_wo_item,
                    v_wo_id,
                    v_scope,
                    v_process_id,
                    JSON_VALUE(v_elem, '$.plannedQuantity'),
                    JSON_VALUE(v_elem, '$.expectedValue'),
                    created_by_val);

            SET v_inserted_id = LAST_INSERT_ID();
            SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', v_inserted_id);

            IF v_scope = 'ORDER' AND v_wo_id IS NOT NULL THEN
                INSERT IGNORE INTO tmp_sched_wo (id) VALUES (v_wo_id);
            ELSEIF v_scope = 'ITEMS' AND v_wo_item IS NOT NULL THEN
                INSERT IGNORE INTO tmp_sched_woi (id) VALUES (v_wo_item);
            END IF;

            SET i = i + 1;
        END WHILE;

    UPDATE tbl_work_order_item woi
        JOIN tmp_sched_woi t ON t.id = woi.id
    SET woi.item_status = 'Scheduled';

    UPDATE tbl_work_order wo
        JOIN (SELECT DISTINCT woi.work_order_id AS id
              FROM tbl_work_order_item woi
                       JOIN tmp_sched_woi t ON t.id = woi.id) x ON x.id = wo.id
    SET wo.status = 'Scheduled';

    UPDATE tbl_work_order wo
        JOIN tmp_sched_wo t ON t.id = wo.id
    SET wo.status = 'Scheduled';

    COMMIT;

    DROP TEMPORARY TABLE IF EXISTS tmp_sched_wo;
    DROP TEMPORARY TABLE IF EXISTS tmp_sched_woi;

    SELECT total        AS SavedTotalOperationProcess,
           inserted_ids AS InsertedIds;
END;
