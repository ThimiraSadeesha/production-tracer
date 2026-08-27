DROP PROCEDURE IF EXISTS shift_save;
CREATE PROCEDURE shift_save(
    IN shift_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_payload JSON;
    DECLARE v_name VARCHAR(255);
    DECLARE v_start_date DATE;
    DECLARE v_start_time TIME;
    DECLARE v_end_date DATE;
    DECLARE v_end_time TIME;
    DECLARE v_department_id BIGINT;
    DECLARE v_status VARCHAR(100);

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    SET v_payload = shift_val;
    IF v_payload IS NULL OR JSON_TYPE(v_payload) = 'NULL' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Shift payload cannot be empty.';
    END IF;
    IF JSON_TYPE(v_payload) = 'OBJECT' THEN
        SET v_payload = JSON_ARRAY(v_payload);
    END IF;
    IF JSON_TYPE(v_payload) <> 'ARRAY' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Shift payload must be a JSON array.';
    END IF;

    SET total = IFNULL(JSON_LENGTH(v_payload), 0);

    START TRANSACTION;

    WHILE i < total
        DO
            SET v_name = COALESCE(
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].shift'))), 'null'), ''),
                    NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].shiftName'))), 'null'), '')
                         );
            SET v_start_date = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].startDate'))),
                                             'null'), '');
            SET v_start_time = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].startTime'))),
                                             'null'), '');
            SET v_end_date = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].endDate'))),
                                           'null'), '');
            SET v_end_time = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].endTime'))),
                                           'null'), '');
            SET v_department_id = CAST(NULLIF(NULLIF(
                    JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].departmentId'))),
                    'null'), '') AS SIGNED);
            SET v_status = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].status'))),
                                         'null'), '');

            IF v_name IS NULL THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'Shift name is required.';
            END IF;

            INSERT INTO tbl_shift (shift,
                                   start_date,
                                   start_time,
                                   end_date,
                                   end_time,
                                   department_id,
                                   status,
                                   created_by)
            VALUES (v_name,
                    v_start_date,
                    v_start_time,
                    v_end_date,
                    v_end_time,
                    v_department_id,
                    v_status,
                    created_by_val);

            SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', LAST_INSERT_ID());
            SET i = i + 1;
        END WHILE;

    COMMIT;
    SELECT total        AS savedTotalShift,
           inserted_ids AS insertedIds;
END;
