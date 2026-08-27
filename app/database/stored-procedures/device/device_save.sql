DROP PROCEDURE IF EXISTS device_save;
CREATE PROCEDURE device_save(
    IN device_val JSON,
    IN created_by_val VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE machine_i INT DEFAULT 0;
    DECLARE total INT DEFAULT 0;
    DECLARE machine_total INT DEFAULT 0;
    DECLARE saved_total INT DEFAULT 0;
    DECLARE app_version_item VARCHAR(50);
    DECLARE last_seen_at_item LONGTEXT;
    DECLARE status_item VARCHAR(50);
    DECLARE device_type_item VARCHAR(50);
    DECLARE machine_id_item BIGINT;
    DECLARE process_id_item BIGINT;
    DECLARE android_id_item VARCHAR(150);
    DECLARE inserted_ids JSON DEFAULT JSON_ARRAY();

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

    SET total = IFNULL(JSON_LENGTH(device_val), 0);

    WHILE i < total
        DO
            SET app_version_item = JSON_UNQUOTE(JSON_EXTRACT(device_val, CONCAT('$[', i, '].appVersion')));
            SET last_seen_at_item = JSON_UNQUOTE(JSON_EXTRACT(device_val, CONCAT('$[', i, '].lastSeenAt')));
            SET status_item = JSON_UNQUOTE(JSON_EXTRACT(device_val, CONCAT('$[', i, '].status')));
            SET device_type_item = COALESCE(
                    NULLIF(JSON_UNQUOTE(JSON_EXTRACT(device_val, CONCAT('$[', i, '].deviceType'))), 'null'),
                    'Terminal');
            SET machine_id_item = JSON_EXTRACT(device_val, CONCAT('$[', i, '].machineId'));
            SET process_id_item = JSON_EXTRACT(device_val, CONCAT('$[', i, '].processId'));
            SET android_id_item = JSON_UNQUOTE(JSON_EXTRACT(device_val, CONCAT('$[', i, '].androidId')));
            SET machine_total = IFNULL(JSON_LENGTH(JSON_EXTRACT(device_val, CONCAT('$[', i, '].machineIds'))), 0);

            IF machine_total > 0 THEN
                SET machine_i = 0;
                WHILE machine_i < machine_total
                    DO
                        SET machine_id_item = JSON_EXTRACT(
                                device_val,
                                CONCAT('$[', i, '].machineIds[', machine_i, ']')
                                              );
                        INSERT INTO tbl_device (app_version,
                                                last_seen_at,
                                                status,
                                                device_type,
                                                machine_id,
                                                process_id,
                                                android_id,
                                                created_by)
                        VALUES (app_version_item,
                                last_seen_at_item,
                                status_item,
                                device_type_item,
                                machine_id_item,
                                process_id_item,
                                android_id_item,
                                created_by_val);
                        SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', LAST_INSERT_ID());
                        SET saved_total = saved_total + 1;
                        SET machine_i = machine_i + 1;
                    END WHILE;
            ELSE
                INSERT INTO tbl_device (app_version,
                                        last_seen_at,
                                        status,
                                        device_type,
                                        machine_id,
                                        process_id,
                                        android_id,
                                        created_by)
                VALUES (app_version_item,
                        last_seen_at_item,
                        status_item,
                        device_type_item,
                        machine_id_item,
                        process_id_item,
                        android_id_item,
                        created_by_val);
                SET inserted_ids = JSON_ARRAY_APPEND(inserted_ids, '$', LAST_INSERT_ID());
                SET saved_total = saved_total + 1;
            END IF;

            SET i = i + 1;
        END WHILE;

    SELECT saved_total AS savedTotalDevice, inserted_ids AS insertedIds;
    COMMIT;
END;