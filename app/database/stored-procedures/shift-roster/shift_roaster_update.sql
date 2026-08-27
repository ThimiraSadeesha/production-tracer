DROP PROCEDURE IF EXISTS shift_roaster_update;
CREATE PROCEDURE shift_roaster_update(
    IN roster_id_val BIGINT,
    IN date_val DATETIME,
    IN shift_id_val BIGINT,
    IN machine_id_val BIGINT,
    IN status_val VARCHAR(100),
    IN operator_id_val BIGINT,
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


    UPDATE tbl_shift_roster
    SET date        = IFNULL(date_val, date),
        shift_id    = IFNULL(shift_id_val, shift_id),
        machine_id  = IFNULL(machine_id_val, machine_id),
        status      = IFNULL(status_val, status),
        operator_id = IFNULL(operator_id_val, operator_id),
        updated_by  = updated_by_val,
        updated_at  = NOW()
    WHERE id = roster_id_val;
    COMMIT;
    SELECT roster_id_val AS updatedRoasterId;

END;