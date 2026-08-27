DROP PROCEDURE IF EXISTS operation_process_update;
CREATE PROCEDURE operation_process_update(
    IN p_id BIGINT,
    IN p_sequence BIGINT,
    IN p_plan_date_time DATETIME,
    IN p_planned_quantity DECIMAL(15, 2),
    IN p_expected_value DECIMAL(15, 2),
    IN p_estimate_time DOUBLE,
    IN p_machine_id BIGINT,
    IN p_process_id BIGINT,
    IN p_operation_scope VARCHAR(10),
    IN p_status VARCHAR(100),
    IN p_remark VARCHAR(255),
    IN p_actor VARCHAR(255)
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    START TRANSACTION;

    IF p_machine_id IS NOT NULL AND p_machine_id > 0
        AND NOT EXISTS (SELECT 1 FROM tbl_machine WHERE id = p_machine_id) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'The provided machineId does not exist.';
    END IF;

    UPDATE tbl_operation_process
    SET sequence         = IFNULL(p_sequence, sequence),
        plan_date_time   = IFNULL(p_plan_date_time, plan_date_time),
        planned_quantity = IFNULL(p_planned_quantity, planned_quantity),
        expected_value   = IFNULL(p_expected_value, expected_value),
        estimate_time    = IFNULL(p_estimate_time, estimate_time),
        machine_id       = IFNULL(p_machine_id, machine_id),
        process_id       = IFNULL(p_process_id, process_id),
        operation_scope  = IFNULL(p_operation_scope, operation_scope),
        status           = IFNULL(p_status, status),
        remark           = IFNULL(p_remark, remark),
        updated_by       = p_actor
    WHERE id = p_id;

    COMMIT;
    SELECT p_id AS operationProcessId;
END;
