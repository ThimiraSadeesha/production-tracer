DROP PROCEDURE IF EXISTS machine_update;
CREATE PROCEDURE machine_update(
    IN machine_id_val BIGINT,
    IN machine_code_val VARCHAR(100),
    IN machine_name_val VARCHAR(255),
    IN description_val VARCHAR(255),
    IN capabilities_val VARCHAR(255),
    IN status_val VARCHAR(255),
    IN hourly_output_val DOUBLE,
    IN make_ready_time_val DOUBLE,
    IN machine_type_id_val INT,
    IN unit_of_machine_id_val INT,
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

    UPDATE tbl_machine
    SET machine_code       = IFNULL(machine_code_val, machine_code),
        machine_name       = IFNULL(machine_name_val, machine_name),
        description        = IFNULL(description_val, description),
        status             = IFNULL(status_val, status),
        hourly_output      = IFNULL(hourly_output_val, hourly_output),
        machine_type_id    = IFNULL(machine_type_id_val, machine_type_id),
        unit_of_machine_id = IFNULL(unit_of_machine_id_val, unit_of_machine_id),
        make_ready_time    = IFNULL(make_ready_time_val, make_ready_time),
        updated_by         = IFNULL(updated_by_val, updated_by)
    WHERE id = machine_id_val;
    COMMIT;

    SELECT machine_id_val AS updatedMachineId;
END;
