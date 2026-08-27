DROP PROCEDURE IF EXISTS process_save;
CREATE PROCEDURE process_save(
    IN process_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total INT;
    DECLARE attr_i INT DEFAULT 0;
    DECLARE attr_total INT DEFAULT 0;
    DECLARE code_item VARCHAR(100);
    DECLARE name_item VARCHAR(255);
    DECLARE status_item VARCHAR(255);
    DECLARE is_required_item TINYINT(1);
    DECLARE department_id_item BIGINT;
    DECLARE process_id_item BIGINT;
    DECLARE attribute_name_item VARCHAR(255);
    DECLARE attribute_required_item TINYINT(1);
    DECLARE inserted_ids TEXT DEFAULT '';

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1
                @sqlstate = RETURNED_SQLSTATE,
                @errno = MYSQL_ERRNO,
                @message_text = MESSAGE_TEXT;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
            ROLLBACK;
        END;

    SET total = JSON_LENGTH(process_val);

    START TRANSACTION;

    WHILE i < total
        DO
            SET code_item = JSON_UNQUOTE(JSON_EXTRACT(process_val, CONCAT('$[', i, '].code')));
            SET name_item = JSON_UNQUOTE(JSON_EXTRACT(process_val, CONCAT('$[', i, '].name')));
            SET status_item = JSON_UNQUOTE(JSON_EXTRACT(process_val, CONCAT('$[', i, '].status')));
            SET is_required_item =
                    CAST(IFNULL(JSON_EXTRACT(process_val, CONCAT('$[', i, '].isRequired')), 0) AS UNSIGNED);
            SET department_id_item = JSON_EXTRACT(process_val, CONCAT('$[', i, '].departmentId'));


            INSERT INTO tbl_process (code, name, status, is_required, department_id, created_by)
            VALUES (code_item, name_item, status_item, is_required_item, department_id_item, created_by_val);

            SET process_id_item = LAST_INSERT_ID();
            SET inserted_ids = CONCAT_WS(',', inserted_ids, process_id_item);

            SET attr_i = 0;
            SET attr_total = IFNULL(JSON_LENGTH(JSON_EXTRACT(process_val, CONCAT('$[', i, '].processAttributes'))), 0);

            WHILE attr_i < attr_total
                DO
                    SET attribute_name_item = JSON_UNQUOTE(
                            JSON_EXTRACT(process_val, CONCAT('$[', i, '].processAttributes[', attr_i, '].attributeName'))
                                               );
                    SET attribute_required_item = CAST(
                            IFNULL(
                                    JSON_EXTRACT(
                                            process_val,
                                            CONCAT('$[', i, '].processAttributes[', attr_i, '].isRequired')
                                    ),
                                    0
                            ) AS UNSIGNED
                                                 );

                    IF attribute_name_item IS NOT NULL AND attribute_name_item <> '' THEN
                        INSERT INTO tbl_process_attributes (process_id, attribute_name, is_required)
                        VALUES (process_id_item, attribute_name_item, attribute_required_item);
                    END IF;

                    SET attr_i = attr_i + 1;
                END WHILE;

            SET i = i + 1;
        END WHILE;
    COMMIT;

    SELECT total AS numberOfProcess;
END;