DROP PROCEDURE IF EXISTS department_save;
CREATE PROCEDURE department_save(
    IN department_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();
    DECLARE v_payload JSON;
    DECLARE v_code VARCHAR(255);
    DECLARE v_name VARCHAR(255);
    DECLARE v_exists INT DEFAULT 0;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    SET v_payload = department_val;
    IF v_payload IS NULL OR JSON_TYPE(v_payload) = 'NULL' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Department payload cannot be empty.';
    END IF;
    IF JSON_TYPE(v_payload) = 'OBJECT' THEN
        SET v_payload = JSON_ARRAY(v_payload);
    END IF;
    IF JSON_TYPE(v_payload) <> 'ARRAY' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Department payload must be a JSON array.';
    END IF;

    SET total = IFNULL(JSON_LENGTH(v_payload), 0);

    START TRANSACTION;

    WHILE i < total
        DO
            SET v_code = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].departmentCode'))),
                                       'null'), '');
            SET v_name = NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(v_payload, CONCAT('$[', i, '].departmentName'))),
                                       'null'), '');

            IF v_code IS NULL OR v_name IS NULL THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'Department code and name are required.';
            END IF;

            SELECT COUNT(*)
            INTO v_exists
            FROM tbl_department
            WHERE department_code = v_code;

            IF v_exists > 0 THEN
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'Department-Code already exists.';
            END IF;

            INSERT INTO tbl_department (department_code,
                                        department_name,
                                        department_status,
                                        created_by,
                                        created_at)
            VALUES (v_code,
                    v_name,
                    'Active',
                    created_by_val,
                    NOW());

            SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', LAST_INSERT_ID());
            SET i = i + 1;
        END WHILE;

    COMMIT;
    SELECT total        AS savedTotalDepartment,
           inserted_ids AS insertedIds;
END;
