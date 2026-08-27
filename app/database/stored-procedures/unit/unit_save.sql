DROP PROCEDURE IF EXISTS unit_save;
CREATE PROCEDURE unit_save(
    IN unit_code_val VARCHAR(255),
    IN unit_name_val VARCHAR(255),
    IN unit_symbol_val VARCHAR(255),
    IN created_by_val VARCHAR(255)
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

    SELECT COUNT(*) INTO @record_count FROM tbl_unit tu WHERE tu.unit_code = unit_code_val;
    IF @record_count > 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Unit-Code is already exist.';
    END IF;
    INSERT INTO tbl_unit (unit_code, unit_name, unit_symbol, unit_status, created_by)
    VALUES (unit_code_val, unit_name_val, unit_symbol_val, 'Active', created_by_val);
    COMMIT;
    SET @InsertedID = LAST_INSERT_ID();
    SELECT @InsertedID as unitId;

END;