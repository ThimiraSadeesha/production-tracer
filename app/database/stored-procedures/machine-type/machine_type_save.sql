DROP PROCEDURE IF EXISTS machine_type_save;
CREATE PROCEDURE machine_type_save(
    IN mt_code VARCHAR(100),
    IN mt_name VARCHAR(255),
    IN mt_status VARCHAR(255),
    IN mt_process_id BIGINT,
    IN machine_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE mt_id BIGINT;
    DECLARE total_m INT;
    DECLARE j INT DEFAULT 0;

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

    INSERT INTO tbl_machine_type (code, name, machine_type_status, process_id, created_by)
    VALUES (mt_code, mt_name, mt_status, mt_process_id, created_by_val);

    SET mt_id = LAST_INSERT_ID();
    SET total_m = JSON_LENGTH(JSON_EXTRACT(machine_val, '$.machines'));

    WHILE j < total_m
        DO
            INSERT INTO tbl_machine (machine_code,
                                     machine_name,
                                     description,
                                     capabilities,
                                     status,
                                     hourly_output,
                                     machine_type_id,
                                     unit_of_machine_id,
                                     make_ready_time,
                                     created_by)
            VALUES (JSON_UNQUOTE(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].machineCode'))),
                    JSON_UNQUOTE(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].machineName'))),
                    JSON_UNQUOTE(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].description'))),
                    JSON_UNQUOTE(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].capabilities'))),
                    JSON_UNQUOTE(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].status'))),
                    CAST(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].hourlyOutput')) AS DOUBLE),
                    mt_id,
                    CAST(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].unit')) AS UNSIGNED),
                    CAST(JSON_EXTRACT(machine_val, CONCAT('$.machines[', j, '].makeReady')) AS DOUBLE),
                    created_by_val);
            SET j = j + 1;
        END WHILE;
    COMMIT;

    SELECT JSON_OBJECT(
                   'machine_type_id', mt_id,
                   'number_of_machines', total_m
           ) AS result;
END;
