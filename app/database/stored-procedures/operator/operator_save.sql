DROP PROCEDURE IF EXISTS operator_save;
CREATE PROCEDURE operator_save(
    IN operator_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_payload JSON;
    DECLARE v_emp_no VARCHAR(255);
    DECLARE v_name VARCHAR(255);
    DECLARE v_section VARCHAR(255);
    DECLARE v_status VARCHAR(100);
    DECLARE v_department_id BIGINT;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    SET v_payload = operator_val;
    IF v_payload IS NULL OR JSON_TYPE(v_payload) = 'NULL' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Operator payload cannot be empty.';
    END IF;
    IF JSON_TYPE(v_payload) = 'OBJECT' THEN
        SET v_payload = JSON_ARRAY(v_payload);
    END IF;
    IF JSON_TYPE(v_payload) <> 'ARRAY' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Operator payload must be a JSON array.';
    END IF;

    SET total = IFNULL(JSON_LENGTH(v_payload), 0);

    START TRANSACTION;

    WHILE i < total
        DO
            SET v_emp_no = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].empNo'))),
                                         'null'), '');
            SET v_name = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].name'))),
                                       'null'), '');
            SET v_section = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].section'))),
                                          'null'), '');
            SET v_status = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].status'))),
                                         'null'), '');
            SET v_department_id = CAST(NULLIF(NULLIF(
                    JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].departmentId'))),
                    'null'), '') AS SIGNED);

            IF v_emp_no IS NULL OR v_name IS NULL THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'Operator emp no and name are required.';
            END IF;

            INSERT INTO tbl_operator (emp_no, name, section, status, department_id, created_by)
            VALUES (v_emp_no, v_name, v_section, v_status, v_department_id, created_by_val);

            SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', LAST_INSERT_ID());
            SET i = i + 1;
        END WHILE;

    COMMIT;
    SELECT total        AS savedTotalOperator,
           inserted_ids AS insertedIds;
END;
