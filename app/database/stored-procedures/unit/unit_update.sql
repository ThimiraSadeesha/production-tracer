DROP PROCEDURE IF EXISTS unit_update;
CREATE PROCEDURE unit_update(
    IN unit_id INT,
    IN unit_code_val VARCHAR(50),
    IN unit_name_val VARCHAR(255),
    IN unit_symbol_val VARCHAR(255),
    IN unit_status_val VARCHAR(255),
    IN updated_by_val VARCHAR(255)
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

    UPDATE tbl_unit
    SET unit_code = IFNULL(unit_code_val, unit_code),
        unit_name = IFNULL(unit_name_val, unit_name),
        unit_symbol = IFNULL(unit_symbol_val, unit_symbol),
        unit_status = IFNULL(unit_status_val, unit_status),
        updated_by = updated_by_val,
        updated_at = NOW()
    WHERE id = unit_id;

    COMMIT;
    SELECT unit_id AS unitId;
END;
