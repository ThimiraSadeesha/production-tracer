DROP PROCEDURE IF EXISTS machine_save;
CREATE PROCEDURE machine_save(
    IN machine_code_val VARCHAR(100),
    IN machine_name_val VARCHAR(255),
    IN description_val VARCHAR(255),
    IN capabilities_val VARCHAR(255),
    IN status_val VARCHAR(255),
    IN hourly_output_val DOUBLE,
    IN make_ready_time_val DOUBLE,
    IN machine_type_id_val INT,
    IN unit_of_machine_id_val INT,
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
    VALUES (machine_code_val,
            machine_name_val,
            description_val,
            capabilities_val,
            status_val,
            hourly_output_val,
            machine_type_id_val,
            unit_of_machine_id_val,
            make_ready_time_val,
            created_by_val);
    COMMIT;

    SELECT LAST_INSERT_ID() AS machineId;
END;
