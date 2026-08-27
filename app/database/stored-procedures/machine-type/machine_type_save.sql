DROP PROCEDURE IF EXISTS machine_type_save;
CREATE PROCEDURE machine_type_save(
    IN machine_type_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE j INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE machine_total INT DEFAULT 0;
    DECLARE saved_total_machines INT DEFAULT 0;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_payload JSON;
    DECLARE v_code VARCHAR(100);
    DECLARE v_name VARCHAR(255);
    DECLARE v_status VARCHAR(255);
    DECLARE v_process_id BIGINT;
    DECLARE v_mt_id BIGINT;
    DECLARE v_machines JSON;
    DECLARE v_machine_code VARCHAR(100);
    DECLARE v_machine_name VARCHAR(255);
    DECLARE v_description VARCHAR(255);
    DECLARE v_capabilities VARCHAR(255);
    DECLARE v_machine_status VARCHAR(255);
    DECLARE v_hourly_output DOUBLE;
    DECLARE v_unit_id INT;
    DECLARE v_make_ready_time DOUBLE;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    SET v_payload = machine_type_val;
    IF v_payload IS NULL OR JSON_TYPE(v_payload) = 'NULL' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Machine type payload cannot be empty.';
    END IF;
    IF JSON_TYPE(v_payload) = 'OBJECT' THEN
        SET v_payload = JSON_ARRAY(v_payload);
    END IF;
    IF JSON_TYPE(v_payload) <> 'ARRAY' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Machine type payload must be a JSON array.';
    END IF;

    SET total = IFNULL(JSON_LENGTH(v_payload), 0);

    START TRANSACTION;

    WHILE i < total
        DO
            SET v_code = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].code'))),
                                       'null'), '');
            SET v_name = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].name'))),
                                       'null'), '');
            SET v_status = COALESCE(
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].status'))), 'null'), ''),
                    'Active'
                           );
            SET v_process_id = CAST(NULLIF(NULLIF(
                    JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].processId'))),
                    'null'), '') AS SIGNED);

            IF v_code IS NULL OR v_name IS NULL THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'Machine type code and name are required.';
            END IF;

            INSERT INTO tbl_machine_type (code, name, machine_type_status, process_id, created_by)
            VALUES (v_code, v_name, v_status, v_process_id, created_by_val);

            SET v_mt_id = LAST_INSERT_ID();
            SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', v_mt_id);

            SET v_machines = JSON_EXTRACT(v_payload, CONCAT('$[', i, '].machines'));
            SET machine_total = IFNULL(JSON_LENGTH(v_machines), 0);
            SET j = 0;

            WHILE j < machine_total
                DO
                    SET v_machine_code = NULLIF(NULLIF(JSON_UNQUOTE(
                            JSON_EXTRACT(v_machines, CONCAT('$[', j, '].machineCode'))), 'null'), '');
                    SET v_machine_name = NULLIF(NULLIF(JSON_UNQUOTE(
                            JSON_EXTRACT(v_machines, CONCAT('$[', j, '].machineName'))), 'null'), '');
                    SET v_description = NULLIF(NULLIF(JSON_UNQUOTE(
                            JSON_EXTRACT(v_machines, CONCAT('$[', j, '].description'))), 'null'), '');
                    SET v_capabilities = NULLIF(NULLIF(JSON_UNQUOTE(
                            JSON_EXTRACT(v_machines, CONCAT('$[', j, '].capabilities'))), 'null'), '');
                    SET v_machine_status = NULLIF(NULLIF(JSON_UNQUOTE(
                            JSON_EXTRACT(v_machines, CONCAT('$[', j, '].status'))), 'null'), '');
                    SET v_hourly_output = CAST(NULLIF(NULLIF(JSON_UNQUOTE(
                            JSON_EXTRACT(v_machines, CONCAT('$[', j, '].hourlyOutput'))), 'null'), '') AS DOUBLE);
                    SET v_unit_id = CAST(COALESCE(
                            NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_machines, CONCAT('$[', j, '].unitOfMachineId'))),
                                          'null'), ''),
                            NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_machines, CONCAT('$[', j, '].unit'))),
                                          'null'), '')
                                         ) AS SIGNED);
                    SET v_make_ready_time = CAST(COALESCE(
                            NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_machines, CONCAT('$[', j, '].makeReadyTime'))),
                                          'null'), ''),
                            NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_machines, CONCAT('$[', j, '].makeReady'))),
                                          'null'), ''),
                            0
                                                 ) AS DOUBLE);

                    INSERT INTO tbl_machine (machine_code,
                                             machine_name,
                                             description,
                                             status,
                                             hourly_output,
                                             machine_type_id,
                                             unit_of_machine_id,
                                             make_ready_time,
                                             created_by)
                    VALUES (v_machine_code,
                            v_machine_name,
                            v_description,
                            v_machine_status,
                            v_hourly_output,
                            v_mt_id,
                            v_unit_id,
                            IFNULL(v_make_ready_time, 0),
                            created_by_val);

                    SET saved_total_machines = saved_total_machines + 1;
                    SET j = j + 1;
                END WHILE;

            SET i = i + 1;
        END WHILE;

    COMMIT;
    SELECT total                AS savedTotalMachineType,
           inserted_ids         AS insertedIds,
           saved_total_machines AS savedTotalMachine;
END;
