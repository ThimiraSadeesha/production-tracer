DROP PROCEDURE IF EXISTS machine_save;
CREATE PROCEDURE machine_save(
    IN machine_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_payload JSON;
    DECLARE v_code VARCHAR(100);
    DECLARE v_name VARCHAR(255);
    DECLARE v_description VARCHAR(255);
    DECLARE v_capabilities VARCHAR(255);
    DECLARE v_status VARCHAR(255);
    DECLARE v_hourly_output DOUBLE;
    DECLARE v_make_ready_time DOUBLE;
    DECLARE v_machine_type_id INT;
    DECLARE v_unit_id INT;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    SET v_payload = machine_val;
    IF v_payload IS NULL OR JSON_TYPE(v_payload) = 'NULL' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Machine payload cannot be empty.';
    END IF;
    IF JSON_TYPE(v_payload) = 'OBJECT' THEN
        SET v_payload = JSON_ARRAY(v_payload);
    END IF;
    IF JSON_TYPE(v_payload) <> 'ARRAY' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Machine payload must be a JSON array.';
    END IF;

    SET total = IFNULL(JSON_LENGTH(v_payload), 0);

    START TRANSACTION;

    WHILE i < total
        DO
            SET v_code = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].machineCode'))),
                                       'null'), '');
            SET v_name = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].machineName'))),
                                       'null'), '');
            SET v_description = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].description'))),
                                              'null'), '');
            SET v_capabilities = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].capabilities'))),
                                               'null'), '');
            SET v_status = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].status'))),
                                         'null'), '');
            SET v_hourly_output = CAST(NULLIF(NULLIF(
                    JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].hourlyOutput'))),
                    'null'), '') AS DOUBLE);
            SET v_make_ready_time = CAST(COALESCE(
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].makeReadyTime'))), 'null'),
                           ''),
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].makeReady'))), 'null'), ''),
                    0
                                         ) AS DOUBLE);
            SET v_machine_type_id = CAST(NULLIF(NULLIF(
                    JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].machineTypeId'))),
                    'null'), '') AS SIGNED);
            SET v_unit_id = CAST(COALESCE(
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].unitOfMachineId'))), 'null'),
                           ''),
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].unit'))), 'null'), '')
                                 ) AS SIGNED);

            INSERT INTO tbl_machine (machine_code,
                                     machine_name,
                                     description,
                                     status,
                                     hourly_output,
                                     machine_type_id,
                                     unit_of_machine_id,
                                     make_ready_time,
                                     created_by)
            VALUES (v_code,
                    v_name,
                    v_description,
                    v_status,
                    v_hourly_output,
                    v_machine_type_id,
                    v_unit_id,
                    IFNULL(v_make_ready_time, 0),
                    created_by_val);

            SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', LAST_INSERT_ID());
            SET i = i + 1;
        END WHILE;

    COMMIT;
    SELECT total        AS savedTotalMachine,
           inserted_ids AS insertedIds;
END;
